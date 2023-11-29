package model

type GoogleAdsAPIResponse struct {
	Status  bool        `json:"status"`
	Type    string      `json:"type"`
	Message string      `json:"message"`
	Log     string      `json:"log"`
	Data    interface{} `json:"data"`
}

type TYPEImplGoogleAdsAPI string

const (
	TYPEImplGoogleAdsAPICustomTargetingKey   TYPEImplGoogleAdsAPI = "impl_google_ads_api_custom_targeting_key"
	TYPEImplGoogleAdsAPICustomTargetingValue TYPEImplGoogleAdsAPI = "impl_google_api_custom_targeting_value"
)
