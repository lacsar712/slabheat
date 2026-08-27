package app

import (
	"context"
	"fmt"
	"time"

	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/model"
)

func (a *App) advanceClock(d time.Duration) {
	if mc, ok := a.clk.(*clock.ManualClock); ok {
		mc.Advance(d)
		time.Sleep(time.Millisecond)
	} else {
		time.Sleep(d)
	}
}

func (a *App) bindGasfuelLoop(holder string, ctx context.Context) context.Context {
	a.mu.Lock()
	if cancel, ok := a.gasfuelLoopCancels[holder]; ok {
		cancel()
	}
	child, cancel := context.WithCancel(ctx)
	a.gasfuelLoopCancels[holder] = cancel
	a.mu.Unlock()
	return child
}

func (a *App) cancelGasfuelLoop(holder string) {
	a.mu.Lock()
	if cancel, ok := a.gasfuelLoopCancels[holder]; ok {
		cancel()
		delete(a.gasfuelLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) cancelAllGasfuelLoops() {
	a.mu.Lock()
	for holder, cancel := range a.gasfuelLoopCancels {
		cancel()
		delete(a.gasfuelLoopCancels, holder)
	}
	a.mu.Unlock()
}

func (a *App) CoalFeedTPH() float64 {
	return a.Snapshot().Burner.GasfuelFlowTPH
}

func (a *App) RunGasfuelRamp(ctx context.Context, holder string, targetTPH float64) error {
	loopCtx := a.bindGasfuelLoop(holder, ctx)
	defer a.cancelGasfuelLoop(holder)
	for {
		if err := loopCtx.Err(); err != nil {
			return fmt.Errorf("%w", model.ErrContextDone)
		}
		snap := a.Snapshot()
		current := snap.Burner.GasfuelFlowTPH
		if current >= targetTPH {
			return nil
		}
		comb := snap.Burner
		comb.GasfuelFlowTPH = current + 1.0
		_ = a.store.UpdateBurner(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.GasfuelFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
}

func (a *App) RunCoalFeed(ctx context.Context, holder string, steps int) error {
	defer a.cancelGasfuelLoop(holder)
	_ = a.bindGasfuelLoop(holder, ctx)
	for i := 0; steps <= 0 || i < steps; i++ {
		snap := a.Snapshot()
		comb := snap.Burner
		comb.GasfuelFlowTPH += 0.5
		_ = a.store.UpdateBurner(a.cfg.UnitID, comb)
		a.telemetry.RecordCoalFeed(comb.GasfuelFlowTPH)
		a.advanceClock(100 * time.Millisecond)
	}
	return nil
}
