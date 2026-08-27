package fsm

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/slabheat/internal/model"
)

type SlabzoneFSM struct {
	mu            sync.RWMutex
	state         model.PlantState
	gasfuelPermissive bool
	preheatComplete  bool
	hooks          *HookChain
}

func NewSlabzoneFSM(unitID string) *SlabzoneFSM {
	_ = unitID
	return &SlabzoneFSM{state: model.StateColdStandby, hooks: NewHookChain()}
}

func (f *SlabzoneFSM) Hooks() *HookChain { return f.hooks }

func (f *SlabzoneFSM) State() model.PlantState {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.state
}

func (f *SlabzoneFSM) SetGasfuelPermissive(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.gasfuelPermissive = ok
}

func (f *SlabzoneFSM) SetPreheatComplete(ok bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.preheatComplete = ok
}

func (f *SlabzoneFSM) GasfuelPermissive() bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return f.gasfuelPermissive
}

func (f *SlabzoneFSM) Dispatch(ctx context.Context, event PlantEvent) (model.PlantState, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	select {
	case <-ctx.Done():
		return f.state, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if event == EvTrip {
		from := f.state
		if f.hooks != nil {
			if err := f.hooks.RunBefore(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		f.state = model.StateTrip
		if f.hooks != nil {
			if err := f.hooks.RunAfter(ctx, from, model.StateTrip, event); err != nil {
				return f.state, err
			}
		}
		return f.state, nil
	}
	next, ok := NextState(f.state, event)
	if !ok {
		return f.state, fmt.Errorf("%s from %s: %w", event, f.state, ErrIllegalTransition)
	}
	if event == EvIgnite && !f.gasfuelPermissive {
		return f.state, fmt.Errorf("%w", model.ErrGasfuelPermissive)
	}
	if event == EvPreheatComplete && !f.preheatComplete {
		return f.state, fmt.Errorf("%w", model.ErrPreheatIncomplete)
	}
	from := f.state
	if f.hooks != nil {
		if err := f.hooks.RunBefore(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	f.state = next
	if f.hooks != nil {
		if err := f.hooks.RunAfter(ctx, from, next, event); err != nil {
			return f.state, err
		}
	}
	return f.state, nil
}

func (f *SlabzoneFSM) ForceState(state model.PlantState) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.state = state
}
