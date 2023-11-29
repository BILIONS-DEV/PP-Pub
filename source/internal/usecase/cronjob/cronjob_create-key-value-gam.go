package cronjob

import (
	"encoding/json"
	"fmt"
	"source/internal/entity/dto"
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo/google-ads-api"
	"strconv"
)

type cronJobKeyValue struct {
	*cronJobUC
}

func newCronJobKeyValue(cronJobUC *cronJobUC) *cronJobKeyValue {
	return &cronJobKeyValue{cronJobUC: cronJobUC}
}

type responseDataCreateKeyValue struct {
	NetworkID int64                     `json:"network_id"`
	KeyID     int64                     `json:"key_id"`
	KeyName   string                    `json:"key_name"`
	KeyGamID  int64                     `json:"key_gam_id"`
	Values    []responseDataCreateValue `json:"values"`
}

type responseDataCreateValue struct {
	ValueID    int64  `json:"value_id"`
	ValueName  string `json:"value_name"`
	ValueGamID int64  `json:"value_gam_id"`
}

func (t *cronJobKeyValue) handler(record *model.CronjobModel) (logData string, logError string) {
	var errs []error
	var responseData []responseDataCreateKeyValue

	defer func() {
		if len(errs) > 0 {
			responseError := dto.Fail(errs...)
			bError, _ := json.Marshal(responseError)
			logError = string(bError)
		}
		if len(responseData) > 0 {
			bData, _ := json.Marshal(responseData)
			logData = string(bData)
		}
	}()

	//=> Get All GAM của admin hệ thống pubpower
	gamNetworks, err := t.repos.GamNetwork.FindAllAdmin()
	if err != nil {
		errs = append(errs, err)
		return
	}

	//=> Tạo Key Value cho tất cả các GAM của admin trong hệ thống
	for _, gamNetwork := range gamNetworks {
		//=> Build Config của GAM để request lên Google API
		config, errs2 := t.buildConfigGAM(gamNetwork)
		if len(errs2) > 0 {
			errs = append(errs, errs2...)
			continue
		}

		for _, keyID := range record.ListObjectID {
			//=> Get Key Record
			keyRecord, err := t.repos.KeyValue.FindByID(keyID)
			if err != nil {
				errs = append(errs, errors.NewWithID(err.Error(), strconv.FormatInt(gamNetwork.NetworkID, 10)+"_"+strconv.FormatInt(keyID, 10), "GET_KEY_RECORD (gamNetworkID_keyID)"))
				continue
			}

			//=> Create Key lên GAM qua Google Ads API
			keyGamID, errs2 := t.createKeyGAM(config, keyRecord)
			if len(errs2) > 0 {
				errs = append(errs, errs2...)
				continue
			}

			//=> từ Key ID tiến hành create Value cho Key đó lên GAM qua Google Ads API
			values, errs2 := t.createValueGAM(config, keyRecord, keyGamID)
			if len(errs2) > 0 {
				errs = append(errs, errs2...)
				continue
			}

			//=> Tạo logData
			if len(values) > 0 {
				responseData = append(responseData, responseDataCreateKeyValue{
					NetworkID: gamNetwork.NetworkID,
					KeyID:     keyID,
					KeyName:   keyRecord.KeyName,
					KeyGamID:  keyGamID,
					Values:    values,
				})
			}
		}
	}
	return
}

func (t *cronJobKeyValue) createKeyGAM(config googleAdsApi.InputConfig, keyRecord *model.KeyModel) (keyID int64, errs []error) {
	//=> Build struct Request Create Key GAM
	inputs := googleAdsApi.InputRequestCreateKey{
		Config:  config,
		KeyName: keyRecord.KeyName,
	}

	// Check Key trong bảng log Key Value GAM đã tạo
	keyValueGAMRecord, err := t.repos.KeyValueGam.FindKeyInGAM(inputs.Config.NetworkID, keyRecord.KeyName)
	if err != nil {
		errs = append(errs, errors.NewWithID(err.Error(), strconv.FormatInt(inputs.Config.NetworkID, 10)+"_"+keyRecord.KeyName, "CHECK_LOG_KEY_GAM", "{networkID}_{keyName}"))
		return
	}

	//=> Nếu key đã được tạo trả ra keyID luôn để xử lý tiếp tạo Value lên GAM
	if keyValueGAMRecord.IsFound() {
		keyID = keyValueGAMRecord.KeyID
		return
	}

	//=> Nếu chưa có log key này trong DB tiến hành request create Key lên Google Ads API
	response, err := t.repos.GoogleAdsAPI.Create(inputs, model.TYPEImplGoogleAdsAPICustomTargetingKey)
	if err != nil {
		errs = append(errs, err)
		return
	}
	keyID = int64(response.Data.(float64))
	return
}

func (t *cronJobKeyValue) createValueGAM(config googleAdsApi.InputConfig, keyRecord *model.KeyModel, keyGamID int64) (responses []responseDataCreateValue, errs []error) {
	//=> Build struct Request create Value GAM
	inputs := googleAdsApi.InputRequestCreateValue{
		Config: config,
		KeyID:  fmt.Sprintf("%v", keyGamID),
	}

	for _, value := range keyRecord.Value {
		//=> Gán value cho request
		inputs.Value = value.Value

		//=> Check Key Value trong GAM từ bảng log Key Value Gam đã tạo
		keyValueGAMRecord, err := t.repos.KeyValueGam.FindKeyValueInGAM(config.NetworkID, keyRecord.KeyName, value.Value)
		if err != nil {
			errs = append(errs, errors.NewWithID(err.Error(), strconv.FormatInt(inputs.Config.NetworkID, 10)+"_"+strconv.FormatInt(value.ID, 10), "GET_LOG_VALUE_GAM"))
			return
		}

		//=> Nếu như Value của Key đã được tạo trong GAM và log vào DB rồi thì bỏ qua
		if keyValueGAMRecord.IsFound() {
			continue
		}

		//=> Nếu chưa được tạo tiến hành request create Value lên Google Ads API
		response, err := t.repos.GoogleAdsAPI.Create(inputs, model.TYPEImplGoogleAdsAPICustomTargetingValue)
		if err != nil {
			errs = append(errs, err)
			return
		}

		//=> parse data lấy id của Value trên GAM
		valueGamID := int64(response.Data.(float64))

		//=> Build response
		responses = append(responses, responseDataCreateValue{
			ValueID:    value.ID,
			ValueName:  value.Value,
			ValueGamID: valueGamID,
		})

		//=> Build model log Key Value GAM vào DB
		keyValueGAMRecord = &model.KeyValueGamModel{
			NetworkID: config.NetworkID,
			Name:      keyRecord.KeyName,
			Value:     value.Value,
			KeyID:     keyGamID,
			ValueID:   valueGamID,
		}

		//=> log lại key value của network vào DB
		if err = t.repos.KeyValueGam.Save(keyValueGAMRecord); err != nil {
			errs = append(errs, errors.NewWithID(err.Error(), strconv.FormatInt(config.NetworkID, 10)+"_"+keyRecord.KeyName+"_"+value.Value, "LOG_KEY_VALUE_GAM"))
			return
		}
	}
	return
}

func (t *cronJobKeyValue) buildConfigGAM(gamNetwork model.GamNetworkModel) (config googleAdsApi.InputConfig, errs []error) {
	//=> Get GAM sign in để lấy refreshToken
	gam, err := t.repos.Gam.FindByID(gamNetwork.GamID)
	if err != nil {
		errs = append(errs, errors.NewWithID(err.Error(), strconv.FormatInt(gamNetwork.NetworkID, 10), "GET_GAM_SIGN_IN"))
		return
	}
	config = googleAdsApi.InputConfig{
		RefreshToken: gam.RefreshToken,
		NetworkID:    gamNetwork.NetworkID,
		NetworkName:  gamNetwork.NetworkName,
	}
	return
}

func (t *cronJobKeyValue) validate(record *model.CronjobModel) (err error) {
	if len(record.ListObjectID) == 0 {
		return errors.New("empty data")
	}
	return
}

func (t *cronJobKeyValue) lock(record *model.CronjobModel) (err error) {
	record.Status = model.StatusCronJobPending
	return t.repos.CronJob.Save(record)
}

func (t *cronJobKeyValue) finish(record *model.CronjobModel) (err error) {
	err = t.repos.CronJob.Save(record)
	return
}
