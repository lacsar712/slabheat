package burner

import (
	"context"
	"fmt"
	"math"

	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/model"
)

type Coordinator struct {
	clk     clock.ProcessClock
	burner  *BurnerController
	airflow *AirflowBalancer
	gasfuel    *GasfuelRegulator
	preheat   *clock.PreheatWindow
	ignition *clock.IgnitionDelayWindow
	warmup  *clock.BurnerWarmupWindow
}

func NewCoordinator(clk clock.ProcessClock) *Coordinator {
	return &Coordinator{
		clk:      clk,
		burner:   NewBurnerController(clk),
		airflow:  NewAirflowBalancer(clk),
		gasfuel:     NewGasfuelRegulator(clk),
		preheat:    clock.NewPreheatWindow(clk),
		ignition: clock.NewIgnitionDelayWindow(clk),
		warmup:   clock.NewBurnerWarmupWindow(clk),
	}
}

func (c *Coordinator) Burner() *BurnerController  { return c.burner }
func (c *Coordinator) Airflow() *AirflowBalancer { return c.airflow }
func (c *Coordinator) Gasfuel() *GasfuelRegulator     { return c.gasfuel }

func (c *Coordinator) StartPreheat(ctx context.Context, snap model.PlantSnapshot) (model.BurnerReading, error) {
	select {
	case <-ctx.Done():
		return snap.Burner, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	out := snap.Burner
	out.BurnerPhase = model.BurnerPreheat
	out.PreheatStartedAt = c.clk.Now()
	out.GasfuelFlowTPH = 0
	out.AirflowTPH = c.airflow.PreheatRate()
	return out, nil
}

func (c *Coordinator) CompletePreheat(snap model.BurnerReading) error {
	return c.preheat.Require(snap.PreheatStartedAt)
}

func (c *Coordinator) Ignite(ctx context.Context, snap model.PlantSnapshot) (model.BurnerReading, error) {
	select {
	case <-ctx.Done():
		return snap.Burner, fmt.Errorf("%w", model.ErrContextDone)
	default:
	}
	if err := c.preheat.Require(snap.Burner.PreheatStartedAt); err != nil {
		return snap.Burner, err
	}
	out := snap.Burner
	out.BurnerPhase = model.BurnerIgnition
	out.IgnitionAt = c.clk.Now()
	out.GasfuelFlowTPH = c.gasfuel.IgnitionRate(snap.Settings)
	out.AirflowTPH = c.airflow.IgnitionRate(snap.Settings)
	out.ReheatfnTempF = 400
	return out, nil
}

func (c *Coordinator) Stabilize(snap model.PlantSnapshot) (model.BurnerReading, error) {
	if err := c.ignition.Require(snap.Burner.IgnitionAt); err != nil {
		return snap.Burner, err
	}
	out := snap.Burner
	out.BurnerPhase = model.BurnerStable
	out.GasfuelFlowTPH = snap.Settings.GasfuelFlowTPH * 0.5
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.ReheatfnTempF = c.burner.EstimateReheatfnTemp(out)
	return out, nil
}

func (c *Coordinator) RampToLoad(snap model.PlantSnapshot, loadPct float64) model.BurnerReading {
	out := snap.Burner
	out.GasfuelFlowTPH = snap.Settings.GasfuelFlowTPH * loadPct
	out.AirflowTPH = c.airflow.Compute(snap)
	out.ExcessO2Pct = c.airflow.ExcessO2(out)
	out.ReheatfnTempF = c.burner.EstimateReheatfnTemp(out)
	return out
}

func (c *Coordinator) Trip(snap model.BurnerReading) model.BurnerReading {
	out := snap
	out.BurnerPhase = model.BurnerTrip
	out.GasfuelFlowTPH = 0
	out.ReheatfnTempF = math.Max(200, out.ReheatfnTempF*0.5)
	return out
}

func (c *Coordinator) WarmupReady(snap model.BurnerReading) bool {
	return c.warmup.Ready(snap.IgnitionAt)
}
