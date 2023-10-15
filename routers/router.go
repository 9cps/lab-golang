package router

import (
	"github.com/9cps/api-go-gin/controllers"

	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(healthCheckController *controllers.HealthCheckController, expensesController *controllers.ExpensesController) *gin.Engine {
	router := gin.Default()
	// add swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "welcome home")
	})

	baseRouter := router.Group("/api/v1")
	healthCheckRouter := baseRouter.Group("/HealthCheck")
	healthCheckRouter.GET("/Api", healthCheckController.HealthCheckAPI)
	healthCheckRouter.GET("/Database", healthCheckController.HealthCheckDB)

	expensesRouter := baseRouter.Group("/Expenses")
	expensesRouter.POST("/CreateExpenses", expensesController.CreateExpenses)
	expensesRouter.POST("/CreateExpensesDetail", expensesController.CreateExpensesDetail)
	expensesRouter.GET("/GetListMoneyCard", expensesController.GetListMoneyCard)
	expensesRouter.POST("/GetListMoneyCardDetail", expensesController.GetListMoneyCardDetail)
	return router
}
