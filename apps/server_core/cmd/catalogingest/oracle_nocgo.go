//go:build !cgo

package main

import (
	"context"
	"database/sql"
	"errors"

	oracleconfig "marketplace-central/apps/server_core/internal/modules/internal_read/adapters/oracle"
)

// openOracleDB is the non-cgo stand-in. It fails closed rather than fabricate
// a connection: a host without a C toolchain cannot speak to Oracle at all,
// and pretending otherwise would be the exact invisible failure this command
// exists to end. See internal/modules/internal_read/adapters/oracle/open_nocgo.go
// for the pattern this mirrors.
func openOracleDB(context.Context, oracleconfig.Config) (*sql.DB, error) {
	return nil, errors.New("catalogingest: oracle driver requires a cgo build with a C toolchain and the Oracle Instant Client; this binary must run inside the environment described in docs/operations/live-oracle-docker.md")
}
