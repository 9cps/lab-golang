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

	_ "github.com/9cps/api-go-gin/docs"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/config"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/handler"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/repository"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/router"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/service"
)

//	@title			Expenses Service API
//	@version		1.0
//	@description	Expenses microservice — manages expense records and details.
//	@host			localhost:8081
//	@BasePath		/api/v1

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345"
func main() {
	config.LoadEnv()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	// Repository
	healthCheckRepository := repository.NewHealthCheckRepository(db)
	expensesRepository := repository.NewExpensesRepository(db)

	// Service
	healthCheckService := service.NewHealthCheckService(healthCheckRepository)
	expensesService := service.NewExpensesService(expensesRepository)

	// Handler
	healthCheckHandler := handler.NewHealthCheckHandler(healthCheckService)
	expensesHandler := handler.NewExpensesHandler(expensesService)

	// Router
	routes := router.NewRouter(healthCheckHandler, expensesHandler)

	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8081"
	}

	server := &http.Server{
		Addr:              serverAddr,
		Handler:           routes,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("expenses-service listening on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("http server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("shutting down expenses-service...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
}
