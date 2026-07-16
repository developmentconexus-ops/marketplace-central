package background

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

type fakePasser struct {
	calls []string
	errOn string
}

func (p *fakePasser) Pass(_ context.Context, installationID string) (bool, error) {
	p.calls = append(p.calls, installationID)
	if installationID == p.errOn {
		return true, errors.New("pass failed")
	}
	return true, nil
}

func runPoller(t *testing.T, p *Poller, ctx context.Context) chan struct{} {
	t.Helper()
	done := make(chan struct{})
	go func() {
		p.Run(ctx)
		close(done)
	}()
	return done
}

func waitClosed(t *testing.T, done chan struct{}, what string) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatalf("%s did not stop", what)
	}
}

func TestPollerRunTicksDrivePassesAndErrorsContinue(t *testing.T) {
	passer := &fakePasser{errOn: "inst-b"}
	installations := func(context.Context) ([]string, error) { return []string{"inst-a", "inst-b", "inst-c"}, nil }
	ticks := make(chan time.Time)
	p := NewPoller(passer, installations, slog.Default(), time.Second, ticks)
	ctx, cancel := context.WithCancel(context.Background())
	done := runPoller(t, p, ctx)

	ticks <- time.Time{}
	ticks <- time.Time{}
	cancel()
	waitClosed(t, done, "poller")

	want := []string{"inst-a", "inst-b", "inst-c", "inst-a", "inst-b", "inst-c"}
	if len(passer.calls) != len(want) {
		t.Fatalf("passes=%v, want %v (inst-b error must not halt the loop)", passer.calls, want)
	}
	for i := range want {
		if passer.calls[i] != want[i] {
			t.Fatalf("passes=%v, want %v", passer.calls, want)
		}
	}
}

func TestPollerRunInstallationListErrorContinuesToNextTick(t *testing.T) {
	passer := &fakePasser{}
	failFirst := true
	installations := func(context.Context) ([]string, error) {
		if failFirst {
			failFirst = false
			return nil, errors.New("list failed")
		}
		return []string{"inst-a"}, nil
	}
	ticks := make(chan time.Time)
	p := NewPoller(passer, installations, slog.Default(), time.Second, ticks)
	ctx, cancel := context.WithCancel(context.Background())
	done := runPoller(t, p, ctx)

	ticks <- time.Time{}
	ticks <- time.Time{}
	cancel()
	waitClosed(t, done, "poller")

	if len(passer.calls) != 1 || passer.calls[0] != "inst-a" {
		t.Fatalf("passes=%v, want one inst-a pass after recovered list error", passer.calls)
	}
}

func TestPollerRunStopsOnCancelWithoutTick(t *testing.T) {
	passer := &fakePasser{}
	installations := func(context.Context) ([]string, error) { return nil, nil }
	ticks := make(chan time.Time)
	p := NewPoller(passer, installations, slog.Default(), time.Second, ticks)
	ctx, cancel := context.WithCancel(context.Background())
	done := runPoller(t, p, ctx)

	cancel()
	waitClosed(t, done, "poller")
	if len(passer.calls) != 0 {
		t.Fatalf("passes=%v, want none", passer.calls)
	}
}
