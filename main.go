package main

import (
	"net/http"
	"os"

	"github.com/9cps/api-go-gin/controllers"
	_ "github.com/9cps/api-go-gin/docs"
	"github.com/9cps/api-go-gin/helper"
	"github.com/9cps/api-go-gin/initializers"
	"github.com/9cps/api-go-gin/repository"
	router "github.com/9cps/api-go-gin/routers"
	services "github.com/9cps/api-go-gin/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//	@title			Swagger Example Golang APIs
//	@version		1.0
//	@description	This is a sample server celler server.

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.basic	BasicAuth
func init() {
	initializers.LoadEnv()
	initializers.ConncetDatabse()
}

var DB *gorm.DB

func main() {
	r := gin.Default()
	// Use CORS middleware with a wildcard to allow all origins for all routes.
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	}))

	// Serve Swagger documentation
	// Repository
	healthCheckRepository := repository.NewHealthCheckRepositoryImpl(DB)
	expensesRepository := repository.NewExpensesRepositoryImpl(DB)

	// Service
	healthCheckServices := services.NewHealthCheckServiceImpl(healthCheckRepository)
	expensesServices := services.NewExpensesServiceImpl(expensesRepository)

	// Controller
	healthCheckController := controllers.NewHealthCheckController(healthCheckServices)
	expensesController := controllers.NewExpensesController(expensesServices)

	// Router
	routes := router.NewRouter(healthCheckController, expensesController)

	// Read the server address from the environment variable
	serverAddr := os.Getenv("SERVER_ADDR")

	// If the SERVER_ADDR environment variable is not set, use a default address
	if serverAddr == "" {
		serverAddr = ":8080" // Set a default address
	}

	server := &http.Server{
		Addr:    serverAddr,
		Handler: routes,
	}

	err := server.ListenAndServe()
	helper.ErrorPanic(err)

	//v1 := r.Group("/api/v1")
	//{
	// HealthCheck := v1.Group("/healthcheck")
	// {
	// 	HealthCheck.GET("/HealthCheckAPI", controllers.HealthCheckAPI)
	// 	HealthCheck.GET("/HealthCheckDB", controllers.HealthCheckDB)
	// }

	// Expense := v1.Group("/expense")
	// {
	// 	Expense.GET("/GetListMoneyCard", controllers.GetListMoneyCard)
	// 	Expense.POST("/CreateExpenses", controllers.CreateExpensesDetail)
	// 	Expense.POST("/CreateExpensesDetail", controllers.CreateExpensesDetail)
	// 	Expense.POST("/GetListMoneyCardDetail", controllers.GetListMoneyCardDetail)
	// }
	//}
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//r.Run()
}
