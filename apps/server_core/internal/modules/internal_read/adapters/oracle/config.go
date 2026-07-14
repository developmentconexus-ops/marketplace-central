package oracle

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	defaultPoolMinSessions  = 1
	defaultPoolMaxSessions  = 12
	defaultPoolIncrement    = 1
	defaultPoolWaitTimeout  = 5 * time.Second
	defaultSessionTimeout   = 5 * time.Minute
	defaultSessionLifetime  = 30 * time.Minute
	defaultConnectTimeout   = 10 * time.Second
	defaultCallTimeout      = 30 * time.Second
	defaultBootstrapTimeout = 30 * time.Second
)

type Config struct {
	Username         string
	Password         string
	ConnectString    string
	LibDir           string
	PoolMinSessions  int
	PoolMaxSessions  int
	PoolIncrement    int
	PoolWaitTimeout  time.Duration
	SessionTimeout   time.Duration
	SessionLifetime  time.Duration
	ConnectTimeout   time.Duration
	CallTimeout      time.Duration
	BootstrapTimeout time.Duration
}

func LoadConfigFromEnv(getenv func(string) string) (Config, error) {
	cfg := Config{
		Username:      strings.TrimSpace(getenv("MPC_ORACLE_USERNAME")),
		Password:      strings.TrimSpace(getenv("MPC_ORACLE_PASSWORD")),
		ConnectString: strings.TrimSpace(getenv("MPC_ORACLE_CONNECT_STRING")),
		LibDir:        strings.TrimSpace(getenv("MPC_ORACLE_LIB_DIR")),
	}
	if cfg.Username == "" {
		return Config{}, errors.New("MPC_ORACLE_USERNAME is required")
	}
	if cfg.Password == "" {
		return Config{}, errors.New("MPC_ORACLE_PASSWORD is required")
	}
	if cfg.ConnectString == "" {
		return Config{}, errors.New("MPC_ORACLE_CONNECT_STRING is required")
	}
	cfg.PoolMinSessions = defaultPoolMinSessions
	cfg.PoolMaxSessions = defaultPoolMaxSessions
	cfg.PoolIncrement = defaultPoolIncrement
	cfg.PoolWaitTimeout = defaultPoolWaitTimeout
	cfg.SessionTimeout = defaultSessionTimeout
	cfg.SessionLifetime = defaultSessionLifetime
	cfg.ConnectTimeout = defaultConnectTimeout
	cfg.CallTimeout = defaultCallTimeout
	cfg.BootstrapTimeout = defaultBootstrapTimeout

	var err error
	if cfg.PoolMinSessions, err = optionalNonNegativeInt(getenv, "MPC_ORACLE_POOL_MIN_SESSIONS", cfg.PoolMinSessions); err != nil {
		return Config{}, err
	}
	if cfg.PoolMaxSessions, err = optionalPositiveInt(getenv, "MPC_ORACLE_POOL_MAX_SESSIONS", cfg.PoolMaxSessions); err != nil {
		return Config{}, err
	}
	if cfg.PoolIncrement, err = optionalPositiveInt(getenv, "MPC_ORACLE_POOL_INCREMENT", cfg.PoolIncrement); err != nil {
		return Config{}, err
	}
	if cfg.PoolWaitTimeout, err = optionalPositiveDuration(getenv, "MPC_ORACLE_POOL_WAIT_TIMEOUT", cfg.PoolWaitTimeout); err != nil {
		return Config{}, err
	}
	if cfg.SessionTimeout, err = optionalPositiveDuration(getenv, "MPC_ORACLE_SESSION_TIMEOUT", cfg.SessionTimeout); err != nil {
		return Config{}, err
	}
	if cfg.SessionLifetime, err = optionalPositiveDuration(getenv, "MPC_ORACLE_SESSION_MAX_LIFETIME", cfg.SessionLifetime); err != nil {
		return Config{}, err
	}
	if cfg.ConnectTimeout, err = optionalPositiveDuration(getenv, "MPC_ORACLE_CONNECT_TIMEOUT", cfg.ConnectTimeout); err != nil {
		return Config{}, err
	}
	if cfg.CallTimeout, err = optionalPositiveDuration(getenv, "MPC_ORACLE_CALL_TIMEOUT", cfg.CallTimeout); err != nil {
		return Config{}, err
	}
	if cfg.BootstrapTimeout, err = optionalPositiveDuration(getenv, "MPC_ORACLE_BOOTSTRAP_TIMEOUT", cfg.BootstrapTimeout); err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if strings.TrimSpace(c.Username) == "" {
		return errors.New("MPC_ORACLE_USERNAME is required")
	}
	if strings.TrimSpace(c.Password) == "" {
		return errors.New("MPC_ORACLE_PASSWORD is required")
	}
	if strings.TrimSpace(c.ConnectString) == "" {
		return errors.New("MPC_ORACLE_CONNECT_STRING is required")
	}
	if c.PoolMaxSessions > 0 && c.PoolMinSessions > c.PoolMaxSessions {
		return fmt.Errorf("MPC_ORACLE_POOL_MIN_SESSIONS cannot exceed MPC_ORACLE_POOL_MAX_SESSIONS")
	}
	if c.PoolMinSessions < 0 {
		return errors.New("MPC_ORACLE_POOL_MIN_SESSIONS must be non-negative")
	}
	if c.PoolMaxSessions <= 0 {
		return errors.New("MPC_ORACLE_POOL_MAX_SESSIONS must be positive")
	}
	if c.PoolIncrement <= 0 || c.PoolIncrement > c.PoolMaxSessions {
		return errors.New("MPC_ORACLE_POOL_INCREMENT must be positive and no greater than MPC_ORACLE_POOL_MAX_SESSIONS")
	}
	for key, value := range map[string]time.Duration{
		"MPC_ORACLE_POOL_WAIT_TIMEOUT":    c.PoolWaitTimeout,
		"MPC_ORACLE_SESSION_TIMEOUT":      c.SessionTimeout,
		"MPC_ORACLE_SESSION_MAX_LIFETIME": c.SessionLifetime,
		"MPC_ORACLE_CONNECT_TIMEOUT":      c.ConnectTimeout,
		"MPC_ORACLE_CALL_TIMEOUT":         c.CallTimeout,
		"MPC_ORACLE_BOOTSTRAP_TIMEOUT":    c.BootstrapTimeout,
	} {
		if value <= 0 {
			return fmt.Errorf("%s must be positive", key)
		}
	}
	if c.ConnectTimeout < time.Second || c.ConnectTimeout%time.Second != 0 {
		return errors.New("MPC_ORACLE_CONNECT_TIMEOUT must be a whole number of seconds and at least 1s")
	}
	return nil
}

func optionalNonNegativeInt(getenv func(string) string, key string, fallback int) (int, error) {
	value := strings.TrimSpace(getenv(key))
	if value == "" {
		return fallback, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("%s must be a non-negative integer", key)
	}
	return parsed, nil
}

func optionalPositiveInt(getenv func(string) string, key string, fallback int) (int, error) {
	value, err := optionalNonNegativeInt(getenv, key, fallback)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return value, nil
}

func optionalPositiveDuration(getenv func(string) string, key string, fallback time.Duration) (time.Duration, error) {
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

func withConnectTimeout(connectString string, timeout time.Duration) string {
	trimmed := strings.TrimSpace(connectString)
	if trimmed == "" || strings.HasPrefix(trimmed, "(") || strings.Contains(strings.ToLower(trimmed), "connect_timeout=") {
		return trimmed
	}
	separator := "?"
	if strings.Contains(trimmed, "?") {
		separator = "&"
	}
	return trimmed + separator + "connect_timeout=" + strconv.FormatInt(int64(timeout/time.Second), 10)
}
