package types

import "time"

type RateLimit struct {
	MaxRequest int
	Duration   time.Duration
	// data user
	UserId   string
	Method   string
	Endpoint string
}
