package router

import (
	"net/http"
	"os"

	"github.com/9cps/api-go-gin/services/api-gateway/internal/middleware"
	"github.com/9cps/api-go-gin/services/api-gateway/internal/proxy"

	"github.com/gin-gonic/gin"
)

// NewRouter builds the api-gateway router.
//
// All requests are forwarded to the expenses-service via reverse proxy.
// The gateway is responsible for CORS and JWT authentication only;
// business logic stays in the upstream services.
func NewRouter() *gin.Engine {
	expensesServiceURL := os.Getenv("EXPENSES_SERVICE_URL")
	if expensesServiceURL == "" {
		expensesServiceURL = "http://localhost:8081"
	}

	upstream := proxy.NewReverseProxy(expensesServiceURL)
	upstream.ErrorHandler = proxy.ErrorHandler
	handler := proxy.Handler(upstream)

	router := gin.Default()
	router.Use(middleware.CorsMiddleware())

	// Welcome
	router.GET("", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "api-gateway")
	})

	// Swagger — proxy to expenses-service (no auth required)
	router.GET("/swagger/*any", handler)

	api := router.Group("/api/v1")

	// Health — public, no auth required
	health := api.Group("/health")
	health.GET("", handler)
	health.GET("/database", handler)

	// Expenses — protected by JWT; gateway validates token, then proxies
	expenses := api.Group("/expenses")
	expenses.Use(middleware.AuthMiddleware())
	{
		expenses.POST("", handler)
		expenses.GET("", handler)
		expenses.POST("/details", handler)
		expenses.GET("/details", handler)
		expenses.PUT("/details/:id", handler)
		expenses.DELETE("/details/:id", handler)
	}

	return router
}
