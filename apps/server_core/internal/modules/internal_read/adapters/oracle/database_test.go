package oracle

import (
	"database/sql"
	"testing"

	"marketplace-central/apps/server_core/internal/modules/internal_read/ports"
)

func TestPoolStatsMappingUsesNeutralValue(t *testing.T) {
	got := poolStatsFromDBStats(sql.DBStats{OpenConnections: 7, InUse: 3, WaitCount: 11})
	want := ports.PoolStats{Open: 7, InUse: 3, WaitCount: 11}
	if got != want {
		t.Fatalf("pool stats = %+v, want %+v", got, want)
	}
}
