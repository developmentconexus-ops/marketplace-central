package background

import (
	"context"
	"log/slog"
	"time"
)

type Passer interface {
	Pass(context.Context, string) (bool, error)
}

type InstallationProvider func(context.Context) ([]string, error)

type Poller struct {
	passer        Passer
	installations InstallationProvider
	logger        *slog.Logger
	interval      time.Duration
	ticks         <-chan time.Time
}

func NewPoller(passer Passer, installations InstallationProvider, logger *slog.Logger, interval time.Duration, ticks <-chan time.Time) *Poller {
	if logger == nil {
		logger = slog.Default()
	}
	if interval <= 0 {
		interval = 2 * time.Second
	}
	return &Poller{passer: passer, installations: installations, logger: logger, interval: interval, ticks: ticks}
}

func (p *Poller) Run(ctx context.Context) {
	if p == nil || p.passer == nil || p.installations == nil {
		return
	}
	ticks := p.ticks
	var ticker *time.Ticker
	if ticks == nil {
		ticker = time.NewTicker(p.interval)
		defer ticker.Stop()
		ticks = ticker.C
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticks:
			installations, err := p.installations(ctx)
			if err != nil {
				p.logger.Error("list mutation poller installations", "err", err)
				continue
			}
			for _, installationID := range installations {
				if _, err := p.passer.Pass(ctx, installationID); err != nil {
					p.logger.Error("mutation poller pass failed", "installation_id", installationID, "err", err)
				}
			}
		}
	}
}
