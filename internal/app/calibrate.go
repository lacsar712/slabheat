package app

import (
	"context"
	"fmt"
)

var CalibrateProbe func(ctx context.Context) error

func (a *App) Calibrate(ctx context.Context, holder string) error {
	now := a.clk.Now()
	_, err := a.interlock.Leases().Acquire(a.cfg.UnitID, holder, now)
	if err != nil {
		return err
	}
	releaseHeld := true
	defer func() {
		if releaseHeld {
			a.interlock.Leases().ReleaseHolder(a.cfg.UnitID, holder)
		}
	}()
	if CalibrateProbe != nil {
		if err := CalibrateProbe(ctx); err != nil {
			return fmt.Errorf("calibrate: %w", err)
		}
	}
	releaseHeld = false
	a.journalEvent("calibrate", fmt.Sprintf("{\"holder\":\"%s\"}", holder))
	return nil
}
