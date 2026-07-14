//go:build !cgo

package oracle

import (
	"os"
	"testing"
)

func TestOracleLiveBaseline(t *testing.T) {
	if os.Getenv("MPC_ORACLE_LIVE_TEST") != "1" {
		t.Skip("set MPC_ORACLE_LIVE_TEST=1 to run live Oracle validation")
	}
	t.Skip("live Oracle baseline requires cgo")
}
