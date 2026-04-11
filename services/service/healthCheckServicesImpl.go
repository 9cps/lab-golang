package services

import (
	"log"

	repoInterfaces "github.com/9cps/api-go-gin/repositories/interfaces"
	svcInterfaces "github.com/9cps/api-go-gin/services/interfaces"
)

type HealthCheckServiceImpl struct {
	HealthCheckRepository repoInterfaces.IHealthCheckRepository
}

// NewHealthCheckServiceImpl returns a HealthCheckServices (service layer)
// rather than a repository interface — the previous return type was
// semantically wrong even though it happened to share a method set.
func NewHealthCheckServiceImpl(healthCheckRepository repoInterfaces.IHealthCheckRepository) svcInterfaces.HealthCheckServices {
	return &HealthCheckServiceImpl{
		HealthCheckRepository: healthCheckRepository,
	}
}

func (s *HealthCheckServiceImpl) HealthCheckDB() bool {
	if ok := s.HealthCheckRepository.HealthCheckDB(); !ok {
		log.Println("health check: database connection failed")
		return false
	}
	return true
}
