package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/9cps/api-go-gin/controllers"
	_ "github.com/9cps/api-go-gin/docs"
	"github.com/9cps/api-go-gin/initializers"
	repositories "github.com/9cps/api-go-gin/repositories/repository"
	router "github.com/9cps/api-go-gin/routers"
	services "github.com/9cps/api-go-gin/services/service"
)

//	@title			Swagger Example Golang APIs
//	@version		1.0
//	@description	This is a sample server celler server.

//	@host		localhost:8080
//	@BasePath	/api/v1

// @securityDefinitions.basic	BasicAuth
func main() {
	initializers.LoadEnv()

	db, err := initializers.ConnectDatabase()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	// Repository
	healthCheckRepository := repositories.NewHealthCheckRepositoryImpl(db)
	expensesRepository := repositories.NewExpensesRepositoryImpl(db)

	// Service
	healthCheckServices := services.NewHealthCheckServiceImpl(healthCheckRepository)
	expensesServices := services.NewExpensesServiceImpl(expensesRepository)

	// Controller
	healthCheckController := controllers.NewHealthCheckController(healthCheckServices)
	expensesController := controllers.NewExpensesController(expensesServices)

	// Router
	routes := router.NewRouter(healthCheckController, expensesController)

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Run the server in a goroutine so we can listen for shutdown signals.
	go func() {
		log.Printf("listening on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	// Wait for SIGINT / SIGTERM for graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	log.Println("server stopped")
}
