package domain

import "errors"

var (
	ErrUnauthorized        = errors.New("provider unauthorized")
	ErrNotFound            = errors.New("provider resource not found")
	ErrRateLimited         = errors.New("provider rate limited")
	ErrProviderUnavailable = errors.New("provider unavailable")
)
