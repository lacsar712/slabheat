package app

import (
	"context"
	"fmt"

	"github.com/lacsar712/slabheat/internal/model"
)

const maxDescalcOpeningPct = 100.0

func (a *App) OpenDescalc(ctx context.Context, holder string, openingPct float64) error {
	_ = holder
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if openingPct >= maxDescalcOpeningPct {
		return fmt.Errorf("descalc: %w", model.ErrDescalcLimit)
	}
	return nil
}

func (a *App) DescalcAfterShutdown(ctx context.Context, openingPct float64) error {
	select {
	case <-ctx.Done():
		return fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	snap := a.Snapshot()
	if snap.State != model.StateTrip && snap.State != model.StateColdStandby {
		return fmt.Errorf("plant not shut down")
	}
	if openingPct >= maxDescalcOpeningPct {
		return fmt.Errorf("descalc: %w", model.ErrDescalcLimit)
	}
	return nil
}
