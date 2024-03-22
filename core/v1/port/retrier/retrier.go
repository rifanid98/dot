package retrier

import (
	"context"
	"dot/core"
	"time"
)

type Effector func(ctx context.Context) (any, *core.CustomError)

//go:generate mockery --name Retrier --filename retrier.go --output ./mocks
type Retrier interface {
	Retry(effector Effector, retries int, delay time.Duration) Effector
}
