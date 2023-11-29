package cronjob

import (
	"source/internal/entity/model"
	"source/internal/errors"
	"source/internal/repo"
)

type UsecaseCronjob interface {
	GetQueues(limit int) (records []model.CronjobModel, err error)
	Handler(record *model.CronjobModel) (errs []error)
}

type cronJobImplInterface interface {
	lock(record *model.CronjobModel) (err error)
	handler(record *model.CronjobModel) (logData string, logError string)
	finish(record *model.CronjobModel) (err error)
	validate(record *model.CronjobModel) (err error)
}

func NewCronJobUC(repos *repo.Repositories) *cronJobUC {
	return &cronJobUC{repos: repos}
}

type cronJobUC struct {
	repos     *repo.Repositories
}

func (t *cronJobUC) GetQueues(limit int) (records []model.CronjobModel, err error) {
	return t.repos.CronJob.FindQueue(limit)
}

func (t *cronJobUC) Handler(record *model.CronjobModel) (errs []error) {
	var impl cronJobImplInterface
	var err error
	var logData, logError string
	if impl, err = t.makeImpl(record.Type); err != nil {
		errs = append(errs, err)
		return
	}
	//=> defer khi return sẽ gọi vào trước khi thoát khỏi hàm
	defer func() {
		//=> Sau khi xử lý xong update thông báo finish cho cronjob
		if err != nil {
			record.Status = model.StatusCronJobError
		} else {
			record.Status = model.StatusCronJobSuccess
		}
		record.Log = logData
		record.Error = logError
		if err = impl.finish(record); err != nil {
			errs = append(errs, err)
			return
		}
	}()

	//=> Lock thông báo pending cho cronjob đang xử lý
	if err = impl.lock(record); err != nil {
		errs = append(errs, err)
		return
	}

	//=> Validate
	if err = impl.validate(record); err != nil {
		errs = append(errs, err)
		return
	}

	//=> Handler
	logData, logError = impl.handler(record)
	if logError != "" {
		err = errors.New("error handler")
		errs = append(errs, err)
		return
	}
	return
}

func (t *cronJobUC) makeImpl(typ model.TYPECronJob) (impl cronJobImplInterface, err error) {
	switch typ {
	case
		model.TYPECronJobCreateKeyValueGAM:
		impl = newCronJobKeyValue(t)

	default:
		return nil, errors.New("Malformed")
	}
	return
}
