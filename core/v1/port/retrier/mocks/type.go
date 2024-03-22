package mocks

import "dot/core/v1/port/retrier"

type RetrierMock struct {
	Retry retrier.Effector
}
