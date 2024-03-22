package util

import (
	"context"
	"time"

	"dot/core"
	"dot/core/v1/port/retrier"
)

var log = NewLogger()

type retrierImpl struct{}

func NewRetrier() *retrierImpl {
	return &retrierImpl{}
}

func (r *retrierImpl) Retry(effector retrier.Effector, retries int, delay time.Duration) retrier.Effector {
	return func(ctx context.Context) (any, *core.CustomError) {
		for r := 0; ; r++ {
			response, err := effector(ctx)
			if err == nil || r >= retries {
				// Return when there is no error or the maximum amount
				// of retries is reached.
				return response, err
			}

			log.Info(ctx, "function call attempt %v failed, retrying next attemp in %v", r+1, delay)

			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return any(nil), &core.CustomError{
					Code:    core.INTERNAL_SERVER_ERROR,
					Message: err.Error(),
				}
			}
		}
	}
}
