package repositories

import (
	"fmt"

	"github.com/9cps/api-go-gin/repositories/interfaces"
	"gorm.io/gorm"
)

type HealthCheckRepositoryImpl struct {
	Db *gorm.DB
}

func NewHealthCheckRepositoryImpl(Db *gorm.DB) interfaces.IHealthCheckRepository {
	return &HealthCheckRepositoryImpl{Db: Db}
}

func (r *HealthCheckRepositoryImpl) HealthCheckDB() bool {
	sqlDB, err := r.Db.DB()
	if err != nil {
		fmt.Printf("Error obtaining sql.DB: %v\n", err)
		return false
	}
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return false
	}
	return true
}
