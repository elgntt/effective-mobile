package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elgntt/effective-mobile/internal/config"
	"github.com/elgntt/effective-mobile/internal/model"
	"github.com/elgntt/effective-mobile/internal/pkg/app_err"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"log/slog"
	"net/http"
	"time"
)

type repository interface {
	AddPeople(ctx context.Context, data model.NewPeople) error
	GetPeople(ctx context.Context, params model.Params) ([]model.People, int, error)
	DeletePeople(ctx context.Context, id int) error
	GetPeopleById(ctx context.Context, id int) (*model.People, error)
	UpdatePeople(ctx context.Context, data model.People) error
}

type Service struct {
	repository
	apiCfg config.APIConfig
	*logrus.Logger
}

func New(repo repository, apiCfg config.APIConfig, log *logrus.Logger) *Service {
	return &Service{
		repo,
		apiCfg,
		log,
	}
}

func (s *Service) AddPeople(ctx context.Context, data model.NewPeopleInfoReq) error {
	s.Logger.Debugln("Getting information via api")
	supplementData, err := s.supplementPeopleInfo(data.Name)
	if err != nil {
		return err
	}

	s.Logger.Debugln("Successfully receiving information from the api")

	return s.repository.AddPeople(ctx, model.NewPeople{
		Name:       data.Name,
		Surname:    data.Surname,
		Patronymic: data.Patronymic,
		Age:        supplementData.Age,
		Gender:     supplementData.Gender,
		CountryID:  supplementData.Country[0].CountryId,
	})
}

func (s *Service) GetPeople(ctx context.Context, params model.Params) ([]model.People, int, error) {
	people, totalCount, err := s.repository.GetPeople(ctx, params)
	if err != nil {
		return nil, 0, err
	}

	return people, totalCount, err
}

func (s *Service) DeletePeople(ctx context.Context, id int) error {
	return s.repository.DeletePeople(ctx, id)
}

func (s *Service) UpdatePeople(ctx context.Context, data model.UpdatePeopleInfoReq) error {
	people, err := s.repository.GetPeopleById(ctx, data.ID)
	if err != nil {
		return err
	}
	if people == nil {
		return app_err.NewBusinessError("People not found")
	}

	if data.Name != nil && people.Name != *data.Name {
		probableData, err := s.supplementPeopleInfo(*data.Name)
		if err != nil {
			return err
		}

		people.Age = probableData.Age
		people.Gender = probableData.Gender
		people.CountryID = probableData.Country[0].CountryId
	}

	people.UpdateFIO(data)

	return s.repository.UpdatePeople(ctx, *people)
}

func (s *Service) supplementPeopleInfo(name string) (model.PeopleDataProbable, error) {
	s.Logger.Debugln("Getting information via API for name:", name)
	urls := []string{
		fmt.Sprintf("%s/?name=%s", s.apiCfg.GetAgeAPI, name),
		fmt.Sprintf("%s/?name=%s", s.apiCfg.GetGenderAPI, name),
		fmt.Sprintf("%s/?name=%s", s.apiCfg.GetCountryAPI, name),
	}

	netClient := http.Client{
		Timeout: 10 * time.Second,
	}

	var (
		additionalData model.PeopleDataProbable
		res            *http.Response
		err            error
		eg             errgroup.Group
	)

	defer func() {
		err = res.Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	for _, url := range urls {
		url := url
		eg.Go(func() error {
			res, err = netClient.Get(url)
			if err != nil {
				return fmt.Errorf("making HTTP request: %w", err)
			}

			err = json.NewDecoder(res.Body).Decode(&additionalData)
			if err != nil {
				return fmt.Errorf("decoding JSON response: %w", err)
			}

			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return model.PeopleDataProbable{}, err
	}

	s.Logger.Debugln("Successfully received information from the API")

	return additionalData, nil
}
