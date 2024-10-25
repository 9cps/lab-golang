package interfaces

type IHealthCheckRepository interface {
	HealthCheckDB() bool
}
