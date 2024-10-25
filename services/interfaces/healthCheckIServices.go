package interfaces

type HealthCheckServices interface {
	HealthCheckDB() bool
}
