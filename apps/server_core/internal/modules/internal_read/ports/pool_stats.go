package ports

// PoolStats is the neutral pool snapshot shared between the Oracle adapter
// and observability without exposing database/sql outside the adapter.
type PoolStats struct {
	Open      int
	InUse     int
	WaitCount int
}
