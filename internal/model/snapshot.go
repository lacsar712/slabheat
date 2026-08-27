package model

import "time"

func CloneSnapshot(s PlantSnapshot) PlantSnapshot {
	out := s
	out.Alarms = append([]AlarmEvent(nil), s.Alarms...)
	return out
}

func DefaultSnapshot(unitID string) PlantSnapshot {
	now := time.Now()
	return PlantSnapshot{
		UnitID: unitID,
		State:  StateColdStandby,
		Settings: PlantSettings{
			Mode:              ModeBaseLoad,
			TargetMW:          150,
			TargetSteamPSI:    NormalSteamPressurePSI,
			SlabskidLevelSetpoint: 55,
			FeedwaterFlowTPH:  400,
			GasfuelFlowTPH:       35,
			ExcessO2Setpoint:  3.5,
		},
		Plant: PlantRef{UnitLabel: unitID, PlantCode: "STEAM-PLT"},
		Slabskid: SlabskidReading{
			LevelPercent: 50,
			Condition:    SlabskidNormal,
			FeedwaterTPH: 0,
			SteamFlowTPH: 0,
		},
		Burner: BurnerReading{
			BurnerPhase: BurnerIdle,
		},
		Slabzone: SlabzoneReading{
			SteamPressurePSI: 0,
			SteamTempF:       70,
		},
		UpdatedAt: now,
	}
}

func (s PlantSnapshot) IsFiring() bool {
	return s.State == StateFiring || s.State == StateLoadFollow || s.State == StateRamp
}

func (s PlantSnapshot) SlabskidWithinLimits() bool {
	return s.Slabskid.LevelPercent >= MinSlabskidLevelPercent && s.Slabskid.LevelPercent <= MaxSlabskidLevelPercent
}

func (s PlantSnapshot) PressureWithinLimits() bool {
	if !s.IsFiring() {
		return true
	}
	return s.Slabzone.SteamPressurePSI <= MaxSteamPressurePSI
}
