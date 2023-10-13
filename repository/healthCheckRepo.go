package repository

type HealthCheckRepository interface {
	HealthCheckDB() bool
}
