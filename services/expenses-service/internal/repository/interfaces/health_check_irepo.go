package interfaces

import "context"

type HealthCheckRepository interface {
	HealthCheckDB(ctx context.Context) (bool, error)
}
