package burner

import (
	"math"

	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/model"
)

type BurnerController struct {
	clk clock.ProcessClock
}

func NewBurnerController(clk clock.ProcessClock) *BurnerController {
	return &BurnerController{clk: clk}
}

func (b *BurnerController) EstimateReheatfnTemp(reading model.BurnerReading) float64 {
	base := 300.0
	gasfuelHeat := reading.GasfuelFlowTPH * 50
	airCool := reading.AirflowTPH * 2
	return base + gasfuelHeat - airCool
}

func (b *BurnerController) ReheatStable(reading model.BurnerReading) bool {
	if reading.BurnerPhase != model.BurnerStable && reading.BurnerPhase != model.BurnerIgnition {
		return false
	}
	return reading.ReheatfnTempF > 800 && reading.ExcessO2Pct >= model.MinReheatfnO2Percent
}

func (b *BurnerController) TripRequired(reading model.BurnerReading) bool {
	if reading.ExcessO2Pct > model.MaxReheatfnO2Percent*2 {
		return true
	}
	if reading.BurnerPhase == model.BurnerTrip {
		return true
	}
	if reading.ReheatfnTempF > 3500 {
		return true
	}
	return false
}

func (b *BurnerController) PhaseLabel(phase model.BurnerPhase) string {
	switch phase {
	case model.BurnerIdle:
		return "Idle"
	case model.BurnerPreheat:
		return "Preheat"
	case model.BurnerIgnition:
		return "Ignition"
	case model.BurnerStable:
		return "Stable Reheat"
	case model.BurnerTrip:
		return "Tripped"
	default:
		return string(phase)
	}
}

func (b *BurnerController) HeatReleaseMW(reading model.BurnerReading) float64 {
	return reading.GasfuelFlowTPH * 12.5
}

func (b *BurnerController) TurndownRatio(settings model.PlantSettings, currentGasfuel float64) float64 {
	if settings.GasfuelFlowTPH <= 0 {
		return 0
	}
	return currentGasfuel / settings.GasfuelFlowTPH
}

func (b *BurnerController) MinStableGasfuel(settings model.PlantSettings) float64 {
	return settings.GasfuelFlowTPH * 0.25
}

func (b *BurnerController) NormalizeGasfuel(flow, max float64) float64 {
	return math.Min(math.Max(flow, 0), max)
}
