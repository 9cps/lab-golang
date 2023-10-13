package services

type HealthCheckServices interface {
	HealthCheckDB() bool
}
