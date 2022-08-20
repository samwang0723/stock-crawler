package services

import "errors"

var (
	ErrWorkerCountInvalid       = errors.New("invalid value for fetch workers")
	ErrRateLimitIntervalInvalid = errors.New("invalid value for rate limit interval")
)
