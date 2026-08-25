package slabskid

import (
	"math"

	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/model"
)

type CarryoverMonitor struct {
	clk clock.ProcessClock
}

func NewCarryoverMonitor(clk clock.ProcessClock) *CarryoverMonitor {
	return &CarryoverMonitor{clk: clk}
}

func (c *CarryoverMonitor) Estimate(slabskid model.SlabskidReading, pressurePSI float64) float64 {
	if slabskid.Condition != model.SlabskidCarry && slabskid.Condition != model.SlabskidSwell {
		return slabskid.CarryoverPPM * 0.9
	}
	base := 50.0
	levelFactor := math.Max(0, slabskid.LevelPercent-70)
	pressureFactor := pressurePSI / 1000
	return base + levelFactor*10 + pressureFactor*5
}

func (c *CarryoverMonitor) AlarmThreshold() float64 { return 500 }

func (c *CarryoverMonitor) TripRequired(ppm float64) bool { return ppm > 1000 }

func (c *CarryoverMonitor) Severity(ppm float64) string {
	switch {
	case ppm > 1000:
		return "critical"
	case ppm > 500:
		return "high"
	case ppm > 200:
		return "medium"
	default:
		return "low"
	}
}

func (c *CarryoverMonitor) RecommendAction(slabskid model.SlabskidReading) string {
	if slabskid.CarryoverPPM > c.AlarmThreshold() {
		return "reduce_load_and_check_separators"
	}
	if slabskid.Condition == model.SlabskidSwell {
		return "hold_feedwater_ramp"
	}
	return "none"
}

func (c *CarryoverMonitor) SeparatorEfficiency(slabskid model.SlabskidReading) float64 {
	eff := 0.98
	if slabskid.LevelPercent > 80 {
		eff -= (slabskid.LevelPercent - 80) * 0.005
	}
	return math.Max(0.5, eff)
}
