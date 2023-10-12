package main

import (
	"net/http"

	"github.com/9cps/api-go-gin/controllers"
	_ "github.com/9cps/api-go-gin/docs"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Swagger Example Golang APIs
//	@version		1.0
//	@description	This is a sample server celler server.

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.basic	BasicAuth
// @Router /healthcheck/* [get]
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
	v1 := r.Group("/api/v1")
	{
		Test := v1.Group("/Test")
		{
			Test.GET("/GetHello", func(c *gin.Context) {
				// Replace this with the actual response data you want to return
				responseData := "Hello, this is a test response."

				c.JSON(http.StatusOK, responseData)
			})
		}

		HealthCheck := v1.Group("/healthcheck")
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
