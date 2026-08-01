package application

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ErrNilTerminalCursor is returned by AssertTerminalCursor when a job body
// reports success (no error) but returns a nil cursor. Per JobFunc's contract
// (see its doc comment), the returned cursor is authoritative: nil does not
// mean "no progress", it means "erase the persisted state" — RecordSuccess
// would write SQL NULL over whatever the entity had synced so far. There is
// no code path in which reporting success with a nil cursor is correct; a job
// with genuinely no new progress must re-return the cursor it was given.
var ErrNilTerminalCursor = errors.New("sync: terminal cursor must be non-nil")

// AssertTerminalCursor validates a job's (cursor, err) result against the
// terminal-cursor contract every registered JobFunc must satisfy, so future
// job authors (M-04/M-06) can prove their body honors it without a database:
//
//  1. A successful cycle (err == nil) must return a non-nil cursor.
//  2. If the cursor carries a "phase" key at all, its value must be one of
//     backfill/incremental/sweep — the vocabulary the scheduler's
//     inferIncremental understands. Anything else is rejected by name rather
//     than silently folded into "not incremental".
//
// A cursor with no "phase" key at all — the existing products job's
// ProductsCursor, for instance — is untouched by rule 2: that job predates
// the phase vocabulary, and this check is opt-in for new phase-aware jobs,
// never retrofitted onto it.
//
// This is a plain function, not a *testing.T-shaped helper, so it is
// importable from any future job's _test.go file without pulling package
// "testing" into this production package. A caller wires it as:
//
//	next, err := myJob(ctx, cursor)
//	if violation := application.AssertTerminalCursor(next, err); violation != nil {
//	    t.Fatal(violation)
//	}
func AssertTerminalCursor(cursor json.RawMessage, jobErr error) error {
	if jobErr != nil {
		// A job that errors is validated by the scheduler's failure path
		// (RecordFailure), not this contract.
		return nil
	}
	if cursor == nil {
		return ErrNilTerminalCursor
	}
	var peek struct {
		Phase string `json:"phase"`
	}
	if err := json.Unmarshal(cursor, &peek); err != nil {
		// Not every valid cursor is a JSON object with a recognizable "phase"
		// key; nothing further to validate here.
		return nil
	}
	if peek.Phase == "" {
		return nil
	}
	switch peek.Phase {
	case "backfill", "incremental", "sweep":
		return nil
	default:
		return fmt.Errorf("sync: cursor phase %q is not one of backfill/incremental/sweep", peek.Phase)
	}
}
