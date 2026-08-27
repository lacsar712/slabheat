package app

import (
	"context"
	"time"
)

func (a *App) RunWarmupBurnScheduler(ctx context.Context, ignitionAt time.Time) error {
	snap, err := a.store.Require(a.cfg.UnitID)
	if err != nil {
		return err
	}
	// Install the burn plan up front so the plan tree shows its child steps
	// while the burner warmup (soak) window is still open. Waiting until the
	// window is satisfied leaves only the root title on the tree during soak,
	// which hides the remaining sub-steps the operator expects to see while
	// the remaining-time indicator is still flashing.
	if err := a.scheduler.InstallBurnPlanCtx(ctx, snap.Settings, "warmup-burn"); err != nil {
		return err
	}
	for !a.warmupWindow.Ready(ignitionAt) {
		if err := ctx.Err(); err != nil {
			return err
		}
		a.advanceClock(100 * time.Millisecond)
	}
	return nil
}

func (a *App) SchedulerItemCount() int {
	return a.scheduler.ItemCount()
}
