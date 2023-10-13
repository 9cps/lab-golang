package services

import (
	"fmt"

	"github.com/9cps/api-go-gin/repository"
)

type HealthCheckServiceImpl struct {
	HealthCheckRepository repository.HealthCheckRepository
}

func NewHealthCheckServiceImpl(healthCheckRepository repository.HealthCheckRepository) HealthCheckServices {
	return &HealthCheckServiceImpl{
		HealthCheckRepository: healthCheckRepository,
	}
}

func (s *HealthCheckServiceImpl) HealthCheckDB() bool {
	statusDb := s.HealthCheckRepository.HealthCheckDB() // Call the function in the repository to open the database return true, false

	if !statusDb {
		fmt.Printf("Error connecting to the database: %v\n", statusDb)
		return false // Return false to indicate a connection failure
	}

	return statusDb // Return true to indicate a successful connection
}
