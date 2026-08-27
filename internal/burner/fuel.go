package burner

import (
	"math"

	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/model"
)

type GasfuelRegulator struct {
	clk clock.ProcessClock
}

func NewGasfuelRegulator(clk clock.ProcessClock) *GasfuelRegulator {
	return &GasfuelRegulator{clk: clk}
}

func (f *GasfuelRegulator) IgnitionRate(settings model.PlantSettings) float64 {
	return settings.GasfuelFlowTPH * 0.08
}

func (f *GasfuelRegulator) ComputeForLoad(settings model.PlantSettings, loadPct float64) float64 {
	loadPct = math.Max(0, math.Min(1, loadPct))
	return settings.GasfuelFlowTPH * loadPct
}

func (f *GasfuelRegulator) Ramp(current, target, maxStep float64) float64 {
	delta := target - current
	if math.Abs(delta) <= maxStep {
		return target
	}
	if delta > 0 {
		return current + maxStep
	}
	return current - maxStep
}

func (f *GasfuelRegulator) BtuPerHour(flowTPH float64) float64 {
	return flowTPH * 19_500_000
}

func (f *GasfuelRegulator) HeatInputMW(flowTPH float64) float64 {
	return flowTPH * 11.6
}

func (f *GasfuelRegulator) ValidatePermissive(settings model.PlantSettings, slabskidOK, preheatOK bool) error {
	if !preheatOK {
		return model.ErrPreheatIncomplete
	}
	if !slabskidOK {
		return model.ErrSlabskidLevelTrip
	}
	if settings.GasfuelFlowTPH <= 0 {
		return model.ErrGasfuelPermissive
	}
	return nil
}

func (f *GasfuelRegulator) MinFlow(settings model.PlantSettings) float64 {
	return settings.GasfuelFlowTPH * 0.2
}

func (f *GasfuelRegulator) MaxFlow(settings model.PlantSettings) float64 {
	return settings.GasfuelFlowTPH * 1.1
}
