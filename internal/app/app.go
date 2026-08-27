package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/lacsar712/slabheat/internal/slabzone"
	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/burner"
	"github.com/lacsar712/slabheat/internal/config"
	"github.com/lacsar712/slabheat/internal/slabskid"
	"github.com/lacsar712/slabheat/internal/fsm"
	"github.com/lacsar712/slabheat/internal/interlock"
	"github.com/lacsar712/slabheat/internal/model"
	"github.com/lacsar712/slabheat/internal/store"
)

type App struct {
	cfg           config.Config
	clk           clock.ProcessClock
	store         *store.PlantStore
	journal       *store.Journal
	fsm           *fsm.SlabzoneFSM
	slabzone        *slabzone.Controller
	burner    *burner.Coordinator
	slabskid          *slabskid.Coordinator
	interlock     *interlock.Interlock
	permissives   *interlock.PermissiveSet
	coordLock     *interlock.CoordinationLock
	scheduler     *clock.Scheduler
	preheatWindow   *clock.PreheatWindow
	warmupWindow  *clock.BurnerWarmupWindow
	telemetry     *Telemetry
	tickCancels    map[string]context.CancelFunc
	gasfuelLoopCancels map[string]context.CancelFunc
	mu             sync.RWMutex
}

func New(cfg config.Config, clk clock.ProcessClock) *App {
	return &App{
		cfg:          cfg,
		clk:          clk,
		store:        store.NewPlantStore(),
		journal:      store.NewJournal(cfg.JournalPath, cfg.JournalCapacity),
		fsm:          fsm.NewSlabzoneFSM(cfg.UnitID),
		slabzone:       slabzone.NewController(clk),
		burner:   burner.NewCoordinator(clk),
		slabskid:         slabskid.NewCoordinator(clk),
		interlock:    interlock.NewInterlock(cfg.LeaseTTL),
		permissives:  interlock.NewPermissiveSet(),
		coordLock:    interlock.NewCoordinationLock(),
		scheduler:    clock.NewScheduler(clk),
		preheatWindow:  clock.NewPreheatWindow(clk),
		warmupWindow: clock.NewBurnerWarmupWindow(clk),
		telemetry:    NewTelemetry(cfg.UnitID),
		tickCancels:     make(map[string]context.CancelFunc),
		gasfuelLoopCancels: make(map[string]context.CancelFunc),
	}
}

func (a *App) Snapshot() model.PlantSnapshot {
	snap, err := a.store.Require(a.cfg.UnitID)
	if err != nil {
		return model.DefaultSnapshot(a.cfg.UnitID)
	}
	return snap
}

func (a *App) Config() config.Config              { return a.cfg }
func (a *App) Clock() clock.ProcessClock          { return a.clk }
func (a *App) FSM() *fsm.SlabzoneFSM                { return a.fsm }
func (a *App) UnitID() string                     { return a.cfg.UnitID }
func (a *App) Store() *store.PlantStore           { return a.store }
func (a *App) Interlock() *interlock.Interlock    { return a.interlock }
func (a *App) Telemetry() TelemetrySnapshot       { return a.telemetry.Snapshot() }
func (a *App) Journal() *store.Journal            { return a.journal }

func (a *App) journalEvent(ev, payload string) {
	_, _ = a.journal.Append(a.cfg.UnitID, ev, payload)
}

func (a *App) syncState(state model.PlantState) {
	_ = a.store.UpdateState(a.cfg.UnitID, state)
}

func (a *App) isFiring(state model.PlantState) bool {
	return state == model.StateFiring || state == model.StateLoadFollow || state == model.StateRamp
}

func (a *App) refreshPermissives(snap model.PlantSnapshot) {
	a.permissives.SetSlabskid(a.slabskid.Level().WithinLimits(snap.Slabskid.LevelPercent))
	a.permissives.SetPressure(a.slabzone.Pressure().WithinTripLimits(snap.Slabzone.SteamPressurePSI, a.isFiring(snap.State)))
	a.permissives.SetBurner(a.burner.Burner().ReheatStable(snap.Burner))
	a.permissives.SetGasfuel(snap.Burner.GasfuelFlowTPH > 0 || snap.State == model.StatePreheat)
	a.permissives.SetIgnition(snap.Burner.BurnerPhase == model.BurnerStable || snap.Burner.BurnerPhase == model.BurnerIgnition)
	a.fsm.SetGasfuelPermissive(a.permissives.GasfuelOK())
	a.fsm.SetPreheatComplete(a.preheatWindow.Ready(snap.Burner.PreheatStartedAt))
}

func (a *App) tickLabel() string {
	return fmt.Sprintf("%s-tick", a.cfg.UnitID)
}
