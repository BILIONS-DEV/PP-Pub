package main

import (
	"source/infrastructure/fakedb"
	"source/internal/repo"
	"source/internal/usecase"
	"testing"
)

func Test_handler_startGetReportTonic(t1 *testing.T) {
	db, err := fakedb.NewMysqlAff()
	if err != nil {
		return
	}
	repos := repo.NewRepositories(&repo.Deps{
		Db:    db,
		Cache: nil,
		Kafka: nil,
	})

	usecases := usecase.NewUseCases(&usecase.Deps{
		Repos:       repos,
		Translation: nil,
	})

	t := &handler{
		UseCases: usecases,
		tonic: tonic{
			token:   "",
			expires: 0,
		},
	}

	t.startGetReportTonic()
}

func Test_handler_startGetReportTaboola(t1 *testing.T) {
	db, err := fakedb.NewMysqlAff()
	if err != nil {
		return
	}
	repos := repo.NewRepositories(&repo.Deps{
		Db:    db,
		Cache: nil,
		Kafka: nil,
	})

	repos.ReportTaboola.Migrate()
	usecases := usecase.NewUseCases(&usecase.Deps{
		Repos:       repos,
		Translation: nil,
	})

	t := &handler{
		UseCases: usecases,
		tonic: tonic{
			token:   "",
			expires: 0,
		},
	}
	t.initTaboola()
	t.startGetReportTaboola()
}

func Test_handler_startGetReportAdsense(t1 *testing.T) {
	db, err := fakedb.NewMysqlAff()
	if err != nil {
		return
	}
	repos := repo.NewRepositories(&repo.Deps{
		Db:    db,
		Cache: nil,
		Kafka: nil,
	})

	usecases := usecase.NewUseCases(&usecase.Deps{
		Repos:       repos,
		Translation: nil,
	})

	t := &handler{
		UseCases: usecases,
		tonic: tonic{
			token:   "",
			expires: 0,
		},
	}

	t.startGetReportAdsense()
}

func Test_handler_updateDataReport(t1 *testing.T) {
	db, err := fakedb.NewMysqlAff()
	if err != nil {
		return
	}
	repos := repo.NewRepositories(&repo.Deps{
		Db:    db,
		Cache: nil,
		Kafka: nil,
	})
	//repos.ReportAff.Migrate()
	usecases := usecase.NewUseCases(&usecase.Deps{
		Repos:       repos,
		Translation: nil,
	})

	t := &handler{
		UseCases: usecases,
		tonic: tonic{
			token:   "",
			expires: 0,
		},
	}

	t.updateDataForReportAff()
}
