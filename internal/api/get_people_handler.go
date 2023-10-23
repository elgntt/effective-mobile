package api

import (
	"context"
	"github.com/elgntt/effective-mobile/internal/model"
	"github.com/elgntt/effective-mobile/internal/pkg/app_err"
	response "github.com/elgntt/effective-mobile/internal/pkg/http"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

func (h Handler) GetPeople(c *gin.Context) {
	ctx := context.Background()

	params, err := parseParams(c.Query("min_age"), c.Query("max_age"), c.Query("gender"),
		c.Query("limit"), c.Query("offset"), c.Query("name"))
	if err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}
	h.logger.Infoln("Getting people information")

	people, totalCount, err := h.service.GetPeople(ctx, params)
	if err != nil {
		response.WriteErrorResponse(c, h.logger, err)
		return
	}

	h.logger.Infoln("People information got successfully")

	c.JSON(http.StatusOK, gin.H{
		"people":     people,
		"totalCount": totalCount,
	})
}

func parsePagination(limitStr, offsetStr string) (int, int, error) {
	if limitStr == "" {
		return 0, 0, app_err.NewBusinessError("the limit parameter is not specified")
	}
	if offsetStr == "" {
		return 0, 0, app_err.NewBusinessError("the offset parameter is not specified")
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return 0, 0, app_err.NewBusinessError("invalid limit parameter")
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, 0, app_err.NewBusinessError("invalid offset parameter")
	}

	return limit, offset, nil
}

func parseParams(minAgeQuery, maxAgeQuery, gender, limit, offset, name string) (model.Params, error) {
	var (
		err    error
		params model.Params
	)
	if minAgeQuery != "" {
		params.MinAge, err = strconv.Atoi(minAgeQuery)
		if err != nil {
			return model.Params{}, app_err.NewBusinessError("Invalid minAge param value")
		}
	}
	if maxAgeQuery != "" {
		params.MaxAge, err = strconv.Atoi(maxAgeQuery)
		if err != nil {
			return model.Params{}, app_err.NewBusinessError("Invalid maxAge param value")
		}
	}

	params.Gender = strings.ToLower(gender)
	params.Name = strings.ToLower(name)
	params.Limit, params.Offset, err = parsePagination(limit, offset)
	if err != nil {
		return model.Params{}, err
	}

	return params, nil
}
