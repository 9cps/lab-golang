package main

import (
	"github.com/9cps/api-go-gin/controllers"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	docs "github.com/9cps/api-go-gin/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /healtcheck/* [get]
// @Router /expense/* [get]

func init() {
	initializers.LoadEnv()
	initializers.ConncetDatabse()
}

func main() {
	r := gin.Default()
	// Use CORS middleware with a wildcard to allow all origins for all routes.
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	}))

	// Serve Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		HealthCheck := v1.Group("/healtcheck")
		{
			HealthCheck.GET("/HealthCheckAPI", controllers.HealthCheckAPI)
			HealthCheck.GET("/HealthCheckDB", controllers.HealthCheckDB)
		}

		// Expenses
		Expense := v1.Group("/expense")
		{
			Expense.GET("/GetListMoneyCard", controllers.GetListMoneyCard)
			Expense.POST("/CreateExpenses", controllers.CreateExpensesDetail)
			Expense.POST("/CreateExpensesDetail", controllers.CreateExpensesDetail)
			Expense.POST("/GetListMoneyCardDetail", controllers.GetListMoneyCardDetail)
		}
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run()
}
