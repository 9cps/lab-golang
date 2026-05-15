package repository

import (
	"context"
	"log/slog"

	"github.com/9cps/api-go-gin/services/expenses-service/internal/repository/interfaces"
	"gorm.io/gorm"
)

type healthCheckRepository struct {
	db *gorm.DB
}

func NewHealthCheckRepository(db *gorm.DB) interfaces.HealthCheckRepository {
	return &healthCheckRepository{db: db}
}

func (r *healthCheckRepository) HealthCheckDB(ctx context.Context) (bool, error) {
	sqlDB, err := r.db.DB()
	if err != nil {
		slog.ErrorContext(ctx, "HealthCheckDB: obtain sql.DB failed", "error", err)
		return false, err
	}
	if err := sqlDB.PingContext(ctx); err != nil {
		slog.ErrorContext(ctx, "HealthCheckDB: ping failed", "error", err)
		return false, err
	}
	return true, nil
}
