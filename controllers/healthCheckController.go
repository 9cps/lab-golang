package controllers

import (
	"net/http"
	"time"

	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	services "github.com/9cps/api-go-gin/services"
	"github.com/gin-gonic/gin"
)

// HealthCheckAPI godoc
//
//	@Summary	Show status api
//	@Tags		healthcheck
//	@Accept		json
//	@Produce	json
//
// @Success      200  {object}  res_dtos.DefaultResponse
//
//	@Router		/healthcheck/HealthCheckAPI [get]
func HealthCheckAPI(c *gin.Context) {
	response := res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "APIs works normally.",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheckDB godoc
//
//	@Summary	Show status database
//	@Tags		healthcheck
//	@Accept		json
//	@Produce	json
//
// @Success      200  {object}  res_dtos.DefaultResponse
//
//	@Router		/healthcheck/HealthCheckDB [get]
func HealthCheckDB(c *gin.Context) {
	db := services.HealthCheckDB()
	var result string

	if db {
		result = "Database connection success."
	} else {
		result = "Database connection failed."
	}

	response := res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: result,
		Date:    time.Now().Format("02/01/2006 15:04:05"),
	}

	c.JSON(200, response)
}
