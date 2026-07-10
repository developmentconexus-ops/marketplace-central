package oracle

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Config struct {
	Username        string
	Password        string
	ConnectString   string
	LibDir          string
	PoolMinSessions int
	PoolMaxSessions int
	SessionTimeout  time.Duration
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
	cfg.PoolMinSessions = 1
	cfg.PoolMaxSessions = 4
	cfg.SessionTimeout = 5 * time.Minute
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
	return nil
}
