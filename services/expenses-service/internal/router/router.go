package router

import (
	"net/http"

	_ "github.com/9cps/api-go-gin/docs"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/handler"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(
	healthCheckHandler handler.HealthCheckHandler,
	expensesHandler handler.ExpensesHandler,
) *gin.Engine {
	r := gin.Default()

	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger/*any", func(ctx *gin.Context) {
		if ctx.Param("any") == "/" {
			ctx.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
			return
		}
		swaggerHandler(ctx)
	})
	r.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "expenses-service")
	})

	api := r.Group("/api/v1")

	// Health check
	health := api.Group("/health")
	health.GET("", healthCheckHandler.HealthCheckAPI)
	health.GET("/database", healthCheckHandler.HealthCheckDB)

	// Expenses — authentication is handled by the api-gateway
	expenses := api.Group("/expenses")
	{
		expenses.GET("", expensesHandler.GetListMoneyCard)
		expenses.PUT("", expensesHandler.CreateExpenses)
		expenses.POST("/details", expensesHandler.GetListMoneyCardDetail)
		expenses.PUT("/details", expensesHandler.UpsertExpensesDetail)
		expenses.DELETE("/details", expensesHandler.DeleteExpensesDetail)
	}

	return r
}
