//go:build cgo

package oracle

import (
	"context"
	"database/sql"
	"time"

	"marketplace-central/apps/server_core/internal/modules/internal_read/domain"
	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"

	"github.com/godror/godror"
)

type managedDB struct {
	*sql.DB
	callTimeout time.Duration
}

func (db *managedDB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, query, append(args, godror.CallTimeout(db.callTimeout))...)
}

func (db *managedDB) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	return db.DB.QueryRowContext(ctx, query, append(args, godror.CallTimeout(db.callTimeout))...)
}

func (db *managedDB) Stats() ports.PoolStats {
	return poolStatsFromDBStats(db.DB.Stats())
}

func OpenDB(ctx context.Context, cfg Config) (Database, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	params := godror.ConnectionParams{}
	params.Username = cfg.Username
	params.Password = godror.NewPassword(cfg.Password)
	params.ConnectString = withConnectTimeout(cfg.ConnectString, cfg.ConnectTimeout)
	params.LibDir = cfg.LibDir
	params.MinSessions = cfg.PoolMinSessions
	params.MaxSessions = cfg.PoolMaxSessions
	params.SessionIncrement = cfg.PoolIncrement
	params.WaitTimeout = cfg.PoolWaitTimeout
	params.SessionTimeout = cfg.SessionTimeout
	params.MaxLifeTime = cfg.SessionLifetime
	params.StandaloneConnection = godror.Bool(false)

	sqlDB := sql.OpenDB(godror.NewConnector(params))
	sqlDB.SetMaxOpenConns(cfg.PoolMaxSessions)
	sqlDB.SetMaxIdleConns(cfg.PoolMinSessions)
	sqlDB.SetConnMaxIdleTime(cfg.SessionTimeout)
	sqlDB.SetConnMaxLifetime(cfg.SessionLifetime)
	db := &managedDB{DB: sqlDB, callTimeout: cfg.CallTimeout}

	bootstrapCtx, cancel := context.WithTimeout(ctx, cfg.BootstrapTimeout)
	defer cancel()
	if err := db.PingContext(bootstrapCtx); err != nil {
		_ = db.Close()
		return nil, domain.NewReadError(domain.ReadErrorSourceUnavailable, "oracle connection ping failed", safeOracleCause(err))
	}
	return db, nil
}
