package api

import (
	"context"
	"github.com/elgntt/effective-mobile/internal/model"
	"github.com/elgntt/effective-mobile/internal/pkg/app_err"
	response "github.com/elgntt/effective-mobile/internal/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) AddPeople(c *gin.Context) {
	ctx := context.Background()

	data := model.NewPeopleInfoReq{}
	if err := c.BindJSON(&data); err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}

	if err := validateData(data); err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}

	h.logger.Infoln("Adding people")
	err := h.service.AddPeople(ctx, data)
	if err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}
	h.logger.Infoln("Successfully added people")
	c.Status(http.StatusCreated)
}

func validateData(data model.NewPeopleInfoReq) error {
	if data.Name == "" {
		return app_err.NewBusinessError("Invalid name!")
	}
	if data.Surname == "" {
		return app_err.NewBusinessError("Invalid surname!")
	}

	return nil
}
