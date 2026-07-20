// Package application holds the sync scheduler skeleton: a cadence-agnostic loop
// that drives registered jobs and keeps sync_state up to date. It reuses the
// proven ticker shape of integrationsbg.NewFeeSyncScheduler (constructor +
// Start(ctx) goroutine + RunOnce cycle) without reinventing it.
package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"marketplace-central/apps/server_core/internal/modules/sync/domain"
)

// StateStore is the persistence seam the scheduler needs. It is implemented by
// adapters/postgres.SyncStateRepository. The repository is constructed with a
// tenant_id and scopes every statement to it (AC-01); callers pass the remaining
// (installation_id, entity) key.
type StateStore interface {
	// Read returns the current state for (installation, entity). found=false
	// means "never synced" — an honest first run, not an error and not a
	// fabricated cursor.
	Read(ctx context.Context, installationID string, entity domain.Entity) (domain.SyncState, bool, error)
	// RecordSuccess upserts the cursor + the relevant timestamp, clears
	// last_error, and resets consecutive_failures to 0. incremental selects
	// last_incremental_at over last_full_sync_at.
	RecordSuccess(ctx context.Context, installationID string, entity domain.Entity, cursor json.RawMessage, at time.Time, incremental bool) error
	// RecordFailure upserts last_error and increments consecutive_failures by 1
	// atomically (single statement). It never touches the success timestamps.
	RecordFailure(ctx context.Context, installationID string, entity domain.Entity, syncErr domain.SyncError) error
}

// JobFunc runs one sync cycle for an entity. It receives the persisted cursor
// (nil on first run) and returns the next cursor to persist. Returning an error
// records last_error + a failure and never brings the loop down. Downstream job
// bodies (M-03/M-04) MUST keep the returned error message generic — it becomes
// sync_state.last_error and must not embed provider credentials or payloads.
type JobFunc func(ctx context.Context, cursor json.RawMessage) (json.RawMessage, error)

// Scheduler drives every registered job on a fixed tick and reconciles
// sync_state. Cadence-agnostic (D6): the tick interval is injected by the caller;
// no business cadence ("daily"/cron) is hardcoded here.
type Scheduler struct {
	store          StateStore
	installationID string
	interval       time.Duration
	now            func() time.Time

	mu   sync.Mutex
	jobs []registeredJob
}

type registeredJob struct {
	entity domain.Entity
	fn     JobFunc
}

// NewScheduler builds a scheduler bound to one installation scope. now may be nil
// (defaults to time.Now) — it exists so tests can inject a deterministic clock.
func NewScheduler(store StateStore, installationID string, interval time.Duration, now func() time.Time) *Scheduler {
	if now == nil {
		now = time.Now
	}
	return &Scheduler{
		store:          store,
		installationID: installationID,
		interval:       interval,
		now:            now,
	}
}

// RegisterJob is the seam by which a future entity plugs a sync body into the
// loop without editing the loop itself (F-02). One registration per entity: a
// duplicate is rejected rather than silently overwritten (first version assumes
// a single job per entity — revisitable). An unknown entity is rejected
// fail-closed (ErrUnknownEntity).
func (s *Scheduler) RegisterJob(entity domain.Entity, fn JobFunc) error {
	if !entity.Valid() {
		return domain.ErrUnknownEntity
	}
	if fn == nil {
		return errors.New("sync: nil job func")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, j := range s.jobs {
		if j.entity == entity {
			return fmt.Errorf("sync: entity %q already registered", entity)
		}
	}
	s.jobs = append(s.jobs, registeredJob{entity: entity, fn: fn})
	return nil
}

// Start runs the loop until ctx is cancelled. A zero/negative interval is a
// no-op (matches the FeeSyncScheduler guard).
func (s *Scheduler) Start(ctx context.Context) {
	if s == nil || s.interval <= 0 {
		return
	}
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.RunOnce(ctx)
		}
	}
}

// RunOnce runs every registered job once. A failing job records last_error and
// increments consecutive_failures but never aborts the remaining jobs — failure
// is isolated per entity.
func (s *Scheduler) RunOnce(ctx context.Context) {
	if s == nil || s.store == nil {
		return
	}
	s.mu.Lock()
	jobs := make([]registeredJob, len(s.jobs))
	copy(jobs, s.jobs)
	s.mu.Unlock()

	for _, j := range jobs {
		if ctx.Err() != nil {
			return
		}
		s.runJob(ctx, j)
	}
}

func (s *Scheduler) runJob(ctx context.Context, j registeredJob) {
	// Fail honest on a read error: a transient failure to read the cursor is NOT
	// a "never synced" first run — fabricating a nil cursor would silently
	// re-sync from scratch (unknown ≠ zero). Skip this entity this cycle; the
	// next tick retries. found=false with no error IS a real first run → nil.
	state, _, err := s.store.Read(ctx, s.installationID, j.entity)
	if err != nil {
		return
	}
	cursor := state.Cursor

	next, jobErr := j.fn(ctx, cursor)
	if jobErr != nil {
		_ = s.store.RecordFailure(ctx, s.installationID, j.entity, domain.SyncError{
			Message: jobErr.Error(),
			At:      s.now().UTC(),
		})
		return
	}
	_ = s.store.RecordSuccess(ctx, s.installationID, j.entity, next, s.now().UTC(), false)
}
