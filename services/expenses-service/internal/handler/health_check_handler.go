package handler

import (
	"net/http"
	"time"

	res "github.com/9cps/api-go-gin/services/expenses-service/internal/dtos/response"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/service/interfaces"
	"github.com/gin-gonic/gin"
)

type HealthCheckHandler interface {
	HealthCheckAPI(ctx *gin.Context)
	HealthCheckDB(ctx *gin.Context)
}

type healthCheckHandler struct {
	healthCheckService interfaces.HealthCheckService
}

func NewHealthCheckHandler(svc interfaces.HealthCheckService) HealthCheckHandler {
	return &healthCheckHandler{healthCheckService: svc}
}

// HealthCheckAPI godoc
//
//	@Summary	Show API status
//	@Tags		HealthCheck
//	@Produce	json
//	@Success	200	{object}	res.DefaultResponse
//	@Router		/health [get]
func (h *healthCheckHandler) HealthCheckAPI(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, res.DefaultResponse{
		Status:  string(res.Success),
		Message: "APIs works normally.",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
	})
}

// HealthCheckDB godoc
//
//	@Summary	Show database connection status
//	@Tags		HealthCheck
//	@Produce	json
//	@Success	200	{object}	res.DefaultResponse
//	@Router		/health/database [get]
func (h *healthCheckHandler) HealthCheckDB(ctx *gin.Context) {
	ok, err := h.healthCheckService.HealthCheckDB(ctx.Request.Context())
	if err != nil || !ok {
		ctx.JSON(http.StatusServiceUnavailable, res.DefaultResponse{
			Status:  string(res.Failed),
			Message: "Database connection failed.",
			Date:    time.Now().Format("02/01/2006 15:04:05"),
		})
		return
	}
	ctx.JSON(http.StatusOK, res.DefaultResponse{
		Status:  string(res.Success),
		Message: "Database connected.",
		Date:    time.Now().Format("02/01/2006 15:04:05"),
	})
}
