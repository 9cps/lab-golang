package controllers

import (
	"net/http"
	"time"

	res_dtos "github.com/9cps/api-go-gin/dtos/response"
	"github.com/9cps/api-go-gin/services/interfaces"
	"github.com/gin-gonic/gin"
)

type HealthCheckController struct {
	healthCheckServices interfaces.HealthCheckServices
}

func NewHealthCheckController(services interfaces.HealthCheckServices) *HealthCheckController {
	return &HealthCheckController{
		healthCheckServices: services,
	}
}

// HealthCheckAPI godoc
//
//	@Summary	Show status api
//	@Tags		HealthCheck
//	@Accept		json
//	@Produce	json
//
//	@Success	200	{object}	res_dtos.DefaultResponse
//
//	@Router		/HealthCheck/Api [GET]
func (c *HealthCheckController) HealthCheckAPI(ctx *gin.Context) {
	response := res_dtos.DefaultResponse{
		Status:  string(res_dtos.Success),
		Message: "APIs works normally.",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
	}

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(http.StatusOK, response)
}

// HealthCheckDB godoc
//
//	@Summary	Show status database
//	@Tags		HealthCheck
//	@Accept		json
//	@Produce	json
//
//	@Success	200	{object}	res_dtos.DefaultResponse
//
//	@Router		/HealthCheck/Database [GET]
func (c *HealthCheckController) HealthCheckDB(ctx *gin.Context) {
	db := c.healthCheckServices.HealthCheckDB()
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

	ctx.Header("Content-Type", "application/json")
	ctx.JSON(200, response)
}
