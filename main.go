package main

import (
	"net/http"
	"os"

	"github.com/9cps/api-go-gin/controllers"
	_ "github.com/9cps/api-go-gin/docs"
	"github.com/9cps/api-go-gin/helper"
	"github.com/9cps/api-go-gin/initializers"
	repositories "github.com/9cps/api-go-gin/repositories/repository"
	router "github.com/9cps/api-go-gin/routers"
	services "github.com/9cps/api-go-gin/services/service"
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
	// Serve Swagger documentation
	// Repository
	healthCheckRepository := repositories.NewHealthCheckRepositoryImpl(DB)
	expensesRepository := repositories.NewExpensesRepositoryImpl(DB)

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
}
