package api

import (
	"context"
	"github.com/elgntt/effective-mobile/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type service interface {
	AddPeople(ctx context.Context, data model.NewPeopleInfoReq) error
	GetPeople(ctx context.Context, params model.Params) ([]model.People, int, error)
	DeletePeople(ctx context.Context, id int) error
	UpdatePeople(ctx context.Context, data model.UpdatePeopleInfoReq) error
}

type Handler struct {
	service service
	logger  *logrus.Logger
}

func New(service service, logger *logrus.Logger) *gin.Engine {
	h := Handler{
		service,
		logger,
	}

	r := gin.New()

	r.POST("/people", h.AddPeople)
	r.GET("/people", h.GetPeople)
	r.DELETE("/people", h.DeletePeople)
	r.PUT("/people", h.UpdatePeople)

	return r
}
