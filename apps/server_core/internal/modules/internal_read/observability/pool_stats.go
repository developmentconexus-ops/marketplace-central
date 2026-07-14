package observability

import (
	"context"
	"database/sql"
	"log/slog"
	"sync"
	"time"
)

type PoolStatsSource interface {
	Stats() sql.DBStats
}

type PoolStatsLoop struct {
	source    PoolStatsSource
	logger    *slog.Logger
	interval  time.Duration
	startOnce sync.Once
	stopOnce  sync.Once
	mu        sync.Mutex
	started   bool
	stop      context.CancelFunc
	done      chan struct{}
}

func NewPoolStatsLoop(source PoolStatsSource, logger *slog.Logger, interval time.Duration) *PoolStatsLoop {
	if logger == nil {
		logger = slog.Default()
	}
	if interval <= 0 {
		interval = defaultPoolStatsInterval
	}
	return &PoolStatsLoop{source: source, logger: logger, interval: interval}
}

func StartPoolStats(ctx context.Context, source PoolStatsSource, logger *slog.Logger, interval time.Duration) *PoolStatsLoop {
	loop := NewPoolStatsLoop(source, logger, interval)
	loop.Start(ctx)
	return loop
}

func (l *PoolStatsLoop) Start(ctx context.Context) {
	if l == nil {
		return
	}
	if ctx == nil {
		ctx = context.Background()
	}
	l.startOnce.Do(func() {
		loopCtx, cancel := context.WithCancel(ctx)
		l.mu.Lock()
		l.started = true
		l.stop = cancel
		l.done = make(chan struct{})
		l.mu.Unlock()

		go func() {
			defer close(l.done)
			if l.source == nil {
				return
			}
			emit := func() {
				stats := l.source.Stats()
				l.logger.Info("pool_stats", "open", stats.OpenConnections, "in_use", stats.InUse, "wait_count", stats.WaitCount)
			}
			emit()
			ticker := time.NewTicker(l.interval)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					emit()
				case <-loopCtx.Done():
					return
				}
			}
		}()
	})
}

func (l *PoolStatsLoop) Stop() {
	if l == nil {
		return
	}
	l.mu.Lock()
	started := l.started
	stop := l.stop
	done := l.done
	l.mu.Unlock()
	if !started {
		return
	}
	l.stopOnce.Do(func() {
		stop()
	})
	<-done
}
