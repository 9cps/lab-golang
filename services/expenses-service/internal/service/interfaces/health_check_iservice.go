package interfaces

import "context"

type HealthCheckService interface {
	HealthCheckDB(ctx context.Context) (bool, error)
}
