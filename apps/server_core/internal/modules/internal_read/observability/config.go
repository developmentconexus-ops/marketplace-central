package observability

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	defaultSlowQueryThreshold = time.Second
	defaultPoolStatsInterval  = 30 * time.Second
)

type Config struct {
	SlowQueryThreshold time.Duration
	PoolStatsInterval  time.Duration
}

func DefaultConfig() Config {
	return Config{
		SlowQueryThreshold: defaultSlowQueryThreshold,
		PoolStatsInterval:  defaultPoolStatsInterval,
	}
}

func LoadConfig(getenv func(string) string) (Config, error) {
	if getenv == nil {
		getenv = os.Getenv
	}
	cfg := DefaultConfig()
	var err error
	if cfg.SlowQueryThreshold, err = optionalDuration(getenv, "MPC_ORACLE_SLOW_QUERY_THRESHOLD", cfg.SlowQueryThreshold); err != nil {
		return Config{}, err
	}
	if cfg.PoolStatsInterval, err = optionalDuration(getenv, "MPC_ORACLE_POOL_STATS_INTERVAL", cfg.PoolStatsInterval); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func optionalDuration(getenv func(string) string, key string, fallback time.Duration) (time.Duration, error) {
	value := strings.TrimSpace(getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := time.ParseDuration(value)
	if err != nil || parsed <= 0 {
		return 0, fmt.Errorf("%s must be a positive Go duration", key)
	}
	return parsed, nil
}
