package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/slabheat/internal/model"
)

func (a *App) WarmupStatus() (ready bool, detail string) {
	snap := a.Snapshot()
	if snap.Burner.PreheatStartedAt.IsZero() {
		return false, "preheat not started"
	}
	if !a.preheatWindow.Ready(snap.Burner.PreheatStartedAt) {
		return false, "preheat window open"
	}
	if !snap.Burner.IgnitionAt.IsZero() && !a.warmupWindow.Ready(snap.Burner.IgnitionAt) {
		return false, "burner warmup window open"
	}
	if !snap.Slabskid.LastSwellAt.IsZero() {
		if err := a.slabskid.RequireSettled(snap.Slabskid); err != nil {
			return false, "slabskid swell settling"
		}
	}
	return true, "ready"
}

func (a *App) WaitWarmup(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("%w", model.ErrContextDone)
		default:
		}
		ready, _ := a.WarmupStatus()
		if ready {
			return nil
		}
	}
}

func (a *App) PreheatRemaining() string {
	snap := a.Snapshot()
	if snap.Burner.PreheatStartedAt.IsZero() {
		return "not started"
	}
	if a.preheatWindow.Ready(snap.Burner.PreheatStartedAt) {
		return "complete"
	}
	return "in progress"
}

func (a *App) BurnerWarmupRemaining() string {
	snap := a.Snapshot()
	if snap.Burner.IgnitionAt.IsZero() {
		return "not ignited"
	}
	if a.warmupWindow.Ready(snap.Burner.IgnitionAt) {
		return "complete"
	}
	return "in progress"
}
