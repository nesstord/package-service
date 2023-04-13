package controllers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"package-service/http/requests"
	"package-service/http/responses"
	"package-service/services"
)

func BoxAggregate(c *gin.Context) {
	request := requests.AggregateRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error(), http.StatusBadRequest))
		return
	}

	response := services.BoxAggregate(c.Request.Context(), request)
	c.JSON(response.ErrorCode, response)
}

func BoxGetBySgtins(c *gin.Context) {
	request := requests.GetBoxesBySgtinsRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error(), http.StatusBadRequest))
		return
	}

	result, err := services.BoxGetBySgtins(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	c.JSON(http.StatusOK, result)
}

func BoxGetByGtin(c *gin.Context) {
	request := requests.GetBoxesByGtinRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, responses.Error(err.Error(), http.StatusBadRequest))
		return
	}

	buf := bytes.Buffer{}
	if err := services.BoxGetByGtin(request, &buf); err != nil {
		c.JSON(http.StatusInternalServerError, responses.Error(err.Error(), http.StatusInternalServerError))
		return
	}

	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Write(buf.Bytes())
}
