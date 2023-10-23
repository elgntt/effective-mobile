package api

import (
	"context"
	"github.com/elgntt/effective-mobile/internal/pkg/app_err"
	response "github.com/elgntt/effective-mobile/internal/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) DeletePeople(c *gin.Context) {
	ctx := context.Background()

	request := struct {
		PeopleId int `json:"peopleId"`
	}{}

	if request.PeopleId < 1 {
		response.WriteErrorResponse(c, h.logger, app_err.NewBusinessError("Invalid people id"))
		return
	}

	if err := c.BindJSON(&request); err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}

	h.logger.Infoln("Deleting people")
	err := h.service.DeletePeople(ctx, request.PeopleId)
	if err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}

	h.logger.Infoln("Deleting people")
	c.Status(http.StatusNoContent)
}
