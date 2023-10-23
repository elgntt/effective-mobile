package http

import (
	"errors"
	"github.com/elgntt/effective-mobile/internal/pkg/app_err"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func WriteErrorResponse(c *gin.Context, logger *logrus.Logger, err error) {
	var bErr app_err.BusinessError

	if errors.As(err, &bErr) {
		errorResponse := ErrorResponse{
			Error: Error{
				Code:    bErr.Code(),
				Message: bErr.Error(),
			},
		}

		logger.Warnln(err.Error())

		c.JSON(http.StatusBadRequest, errorResponse)
	} else {
		errorResponse := ErrorResponse{
			Error: Error{
				Code:    "InternalServerError",
				Message: "Что-то пошло не так, попробуйте еще раз",
			},
		}

		logger.Errorln(err.Error())

		c.JSON(http.StatusInternalServerError, errorResponse)
	}
}
