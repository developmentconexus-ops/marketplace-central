package application

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

// fakeNilJob deliberately violates the terminal-cursor contract: it reports
// success but returns no cursor. This is the must-fail proof for
// AssertTerminalCursor — a contract check that cannot be shown catching a bad
// job is not a check.
func fakeNilJob(context.Context, json.RawMessage) (json.RawMessage, error) {
	return nil, nil
}

// fakeValidJob returns a well-formed, non-nil, phase-tagged cursor — the
// counterpart proof that the check does not also reject a good job.
func fakeValidJob(context.Context, json.RawMessage) (json.RawMessage, error) {
	return json.RawMessage(`{"phase":"incremental"}`), nil
}

func TestAssertTerminalCursorCatchesNilCursor(t *testing.T) {
	next, err := fakeNilJob(context.Background(), nil)
	violation := AssertTerminalCursor(next, err)
	if violation == nil {
		t.Fatal("want the contract check to reject a (nil, nil) job result")
	}
	if !strings.Contains(violation.Error(), "terminal cursor must be non-nil") {
		t.Errorf("violation = %q, want it to contain %q", violation.Error(), "terminal cursor must be non-nil")
	}
}

func TestAssertTerminalCursorPassesAValidCursor(t *testing.T) {
	next, err := fakeValidJob(context.Background(), nil)
	if violation := AssertTerminalCursor(next, err); violation != nil {
		t.Errorf("want a non-nil cursor to pass, got violation: %v", violation)
	}
}

func TestAssertTerminalCursorIgnoresAJobError(t *testing.T) {
	// A job that errors is validated by the scheduler's failure path
	// (RecordFailure), not this contract — an errored job with a nil cursor
	// must not be flagged here.
	if violation := AssertTerminalCursor(nil, errors.New("boom")); violation != nil {
		t.Errorf("want an errored job to bypass the contract, got: %v", violation)
	}
}

func TestAssertTerminalCursorRejectsAnUnrecognizedPhase(t *testing.T) {
	violation := AssertTerminalCursor(json.RawMessage(`{"phase":"weird"}`), nil)
	if violation == nil {
		t.Fatal("want an unrecognized phase to be rejected")
	}
	if !strings.Contains(violation.Error(), "weird") {
		t.Errorf("violation = %q, want it to name the bad phase", violation.Error())
	}
}

func TestAssertTerminalCursorAcceptsEachRecognizedPhase(t *testing.T) {
	for _, phase := range []string{"backfill", "incremental", "sweep"} {
		cursor := json.RawMessage(`{"phase":"` + phase + `"}`)
		if violation := AssertTerminalCursor(cursor, nil); violation != nil {
			t.Errorf("phase %q: want it accepted, got: %v", phase, violation)
		}
	}
}

// TestAssertTerminalCursorAcceptsALegacyCursorWithNoPhaseKey pins that the
// phase-vocabulary check is scoped to new (M-04/M-06) jobs, not retrofitted
// onto the existing products job: its ProductsCursor shape
// (source/processed/completed_at) has no "phase" key at all, and that must
// keep passing with zero behavior change.
func TestAssertTerminalCursorAcceptsALegacyCursorWithNoPhaseKey(t *testing.T) {
	cursor := json.RawMessage(`{"source":"xlsx","processed":1845,"completed_at":"2026-07-24T12:00:00Z"}`)
	if violation := AssertTerminalCursor(cursor, nil); violation != nil {
		t.Errorf("legacy cursor with no phase key must pass, got: %v", violation)
	}
}
