package oracle

import (
	"errors"
	"strings"
)

type Config struct {
	DSN string
}

func LoadConfigFromEnv(getenv func(string) string) (Config, error) {
	dsn := strings.TrimSpace(getenv("SANKHYA_DSN"))
	if dsn == "" {
		return Config{}, errors.New("SANKHYA_DSN is required")
	}
	return Config{DSN: dsn}, nil
}
