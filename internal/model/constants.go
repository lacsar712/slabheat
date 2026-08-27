package model

import "time"

const (
	DefaultLeaseTTL        = 30 * time.Second
	PreheatWindow            = 5 * time.Minute
	IgnitionDelayWindow    = 15 * time.Second
	SlabskidSwellSettleWindow  = 45 * time.Second
	BurnerWarmupWindow = 2 * time.Minute
	FeedwaterRampWindow    = 30 * time.Second
	MaxSlabskidLevelPercent    = 95.0
	MinSlabskidLevelPercent    = 15.0
	TripSlabskidLowPercent     = 10.0
	TripSlabskidHighPercent    = 98.0
	NormalSteamPressurePSI = 1800.0
	MaxSteamPressurePSI    = 2000.0
	MinReheatfnO2Percent    = 2.5
	MaxReheatfnO2Percent    = 6.0
	DefaultJournalCapacity = 512
)
