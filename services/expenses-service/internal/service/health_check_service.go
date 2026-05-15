package service

import (
	"context"
	"fmt"
	"log/slog"

	repoInterfaces "github.com/9cps/api-go-gin/services/expenses-service/internal/repository/interfaces"
	"github.com/9cps/api-go-gin/services/expenses-service/internal/service/interfaces"
)

type healthCheckService struct {
	healthCheckRepo repoInterfaces.HealthCheckRepository
}

func NewHealthCheckService(healthCheckRepo repoInterfaces.HealthCheckRepository) interfaces.HealthCheckService {
	return &healthCheckService{healthCheckRepo: healthCheckRepo}
}

func (s *healthCheckService) HealthCheckDB(ctx context.Context) (bool, error) {
	ok, err := s.healthCheckRepo.HealthCheckDB(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "health check: database connection failed", "error", err)
		return false, fmt.Errorf("health check db: %w", err)
	}
	return ok, nil
}
