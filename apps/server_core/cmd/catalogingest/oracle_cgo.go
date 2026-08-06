//go:build cgo

package main

import (
	"context"
	"database/sql"

	"github.com/godror/godror"

	oracleconfig "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle"
)

// openOracleDB registers the godror driver and opens the connection. This is
// the composition root's job, deliberately: internal/adapters/erp/sankhyaoracle
// never imports godror, so its tests build without cgo. Here, at the edge, the
// driver import is unavoidable and the reason this file — and only this file —
// needs a C toolchain and the Oracle Instant Client to build and link.
func openOracleDB(ctx context.Context, cfg oracleconfig.Config) (*sql.DB, error) {
	params := godror.ConnectionParams{}
	params.Username = cfg.Username
	params.Password = godror.NewPassword(cfg.Password)
	params.ConnectString = cfg.ConnectString
	params.LibDir = cfg.LibDir
	params.MinSessions = cfg.PoolMinSessions
	params.MaxSessions = cfg.PoolMaxSessions
	params.SessionIncrement = cfg.PoolIncrement
	params.WaitTimeout = cfg.PoolWaitTimeout
	params.SessionTimeout = cfg.SessionTimeout
	params.MaxLifeTime = cfg.SessionLifetime
	params.StandaloneConnection = godror.Bool(false)

	db := sql.OpenDB(godror.NewConnector(params))
	db.SetMaxOpenConns(cfg.PoolMaxSessions)
	db.SetMaxIdleConns(cfg.PoolMinSessions)
	db.SetConnMaxIdleTime(cfg.SessionTimeout)
	db.SetConnMaxLifetime(cfg.SessionLifetime)

	bootstrapCtx, cancel := context.WithTimeout(ctx, cfg.BootstrapTimeout)
	defer cancel()
	if err := db.PingContext(bootstrapCtx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
