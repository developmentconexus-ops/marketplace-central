//go:build cgo

package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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
		return nil, safeOracleCause(err)
	}
	return db, nil
}

// safeOracleCause strips raw Oracle/OCI driver text — connection strings, host
// details — out of an error before it can reach main's stderr print.
//
// It is a deliberate copy of the unexported safeOracleCause in
// internal/modules/internal_read/adapters/oracle/reader.go rather than a call
// into it. Exporting the original would have been less code, but internal/modules
// is the frozen legacy tree: it is inventoried, never grown, and widening its
// public surface to serve new code points the dependency the wrong way. This
// file outlives that tree, so it carries its own copy. Keep the two in sync
// until the legacy reader is retired, at which point this one simply remains.
func safeOracleCause(err error) error {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	var coded interface{ Code() int }
	if errors.As(err, &coded) {
		return fmt.Errorf("oracle error code=%d", coded.Code())
	}
	return errors.New("oracle driver error")
}
