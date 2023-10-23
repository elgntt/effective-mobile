package api

import (
	"context"
	"github.com/elgntt/effective-mobile/internal/model"
	response "github.com/elgntt/effective-mobile/internal/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h Handler) UpdatePeople(c *gin.Context) {
	ctx := context.Background()

	data := model.UpdatePeopleInfoReq{}
	if err := c.BindJSON(&data); err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}

	h.logger.Infoln("Updating people information")

	if err := h.service.UpdatePeople(ctx, data); err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}

	h.logger.Infoln("People information updated successfully")
	c.Status(http.StatusOK)
}
