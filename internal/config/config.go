package config

import "time"

type AgentConfig struct {
	ConcurrencyLimit  int
	UserAgent         string
	Timeout           time.Duration
	RetryCount        int
	RequestsPerMinute int
}
