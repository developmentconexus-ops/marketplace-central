package pgdb

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL     string
	DefaultTenantID string
	EncryptionKey   string
}

// LoadPoolConfig reads what opening a pool needs and nothing else. It does not
// read MPC_ENCRYPTION_KEY at all, so a process that never decrypts a credential
// never holds one: the returned Config carries an empty EncryptionKey by
// construction rather than by an unread field happening to stay unused.
//
// Requiring the key everywhere made the opposite true. An offline reprocess,
// which resolves no credential and opens no channel connection, still refused
// to start without a secret in its environment — and an operator's only way
// past that is to export the real key into a process that has no use for it,
// which is exactly the exposure the key exists to avoid.
func LoadPoolConfig() (Config, error) {
	cfg := Config{
		DatabaseURL:     os.Getenv("MC_DATABASE_URL"),
		DefaultTenantID: os.Getenv("MC_DEFAULT_TENANT_ID"),
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("MC_DATABASE_URL is required")
	}
	if cfg.DefaultTenantID == "" {
		cfg.DefaultTenantID = "tenant_default"
	}
	return cfg, nil
}

// LoadConfig is LoadPoolConfig plus the encryption key, for the processes that
// actually decrypt a stored credential. The key stays required there: a
// credential-resolving command that started without one would fail later, at
// the point of use, with a message about decryption rather than about
// configuration.
func LoadConfig() (Config, error) {
	cfg, err := LoadPoolConfig()
	if err != nil {
		return Config{}, err
	}
	cfg.EncryptionKey = os.Getenv("MPC_ENCRYPTION_KEY")
	if cfg.EncryptionKey == "" {
		return Config{}, errors.New("MPC_ENCRYPTION_KEY is required")
	}
	return cfg, nil
}
