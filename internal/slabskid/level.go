package slabskid

import (
	"math"

	"github.com/lacsar712/slabheat/internal/clock"
	"github.com/lacsar712/slabheat/internal/model"
)

type LevelController struct {
	clk clock.ProcessClock
}

func NewLevelController(clk clock.ProcessClock) *LevelController {
	return &LevelController{clk: clk}
}

func (l *LevelController) Compute(snap model.PlantSnapshot, firing bool) (float64, model.SlabskidCondition) {
	level := snap.Slabskid.LevelPercent
	if !firing {
		return level, model.SlabskidNormal
	}
	balance := snap.Slabskid.FeedwaterTPH - snap.Slabskid.SteamFlowTPH
	level += balance * 0.01
	level = math.Max(model.MinSlabskidLevelPercent, math.Min(model.MaxSlabskidLevelPercent, level))
	cond := l.classify(level, snap)
	return level, cond
}

func (l *LevelController) classify(level float64, snap model.PlantSnapshot) model.SlabskidCondition {
	setpoint := snap.Settings.SlabskidLevelSetpoint
	if level > setpoint+15 {
		return model.SlabskidSwell
	}
	if level < setpoint-15 {
		return model.SlabskidShrink
	}
	if snap.Slabzone.SteamPressurePSI > snap.Settings.TargetSteamPSI*0.9 && level > setpoint+5 {
		return model.SlabskidCarry
	}
	return model.SlabskidNormal
}

func (l *LevelController) RecommendFeedwater(snap model.PlantSnapshot, firing bool) float64 {
	if !firing {
		return 0
	}
	err := snap.Settings.SlabskidLevelSetpoint - snap.Slabskid.LevelPercent
	return snap.Settings.FeedwaterFlowTPH + err*3
}

func (l *LevelController) WithinLimits(level float64) bool {
	return level >= model.MinSlabskidLevelPercent && level <= model.MaxSlabskidLevelPercent
}

func (l *LevelController) TripLow(level float64) bool  { return level < model.TripSlabskidLowPercent }
func (l *LevelController) TripHigh(level float64) bool { return level > model.TripSlabskidHighPercent }

func (l *LevelController) LevelError(snap model.PlantSnapshot) float64 {
	return snap.Slabskid.LevelPercent - snap.Settings.SlabskidLevelSetpoint
}

func (l *LevelController) ThreeElementBias(snap model.PlantSnapshot) float64 {
	steam := snap.Slabskid.SteamFlowTPH
	feed := snap.Slabskid.FeedwaterTPH
	levelErr := l.LevelError(snap)
	return feed + (steam-feed)*0.5 + levelErr*2
}
