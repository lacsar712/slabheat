package model

import "errors"

var (
	ErrContextDone      = errors.New("operation cancelled")
	ErrPlantNotFound    = errors.New("plant unit not found")
	ErrLeaseHeld        = errors.New("interlock lease held by another operator")
	ErrLeaseMissing     = errors.New("interlock lease missing or expired")
	ErrGateBlocked      = errors.New("safety gate blocked")
	ErrGasfuelPermissive   = errors.New("gasfuel permissive not satisfied")
	ErrIgnitionBlocked  = errors.New("ignition sequence blocked")
	ErrSlabskidLevelTrip    = errors.New("slabskid level trip condition")
	ErrPressureTrip     = errors.New("steam pressure trip condition")
	ErrBurnerTrip   = errors.New("burner trip condition")
	ErrIllegalState     = errors.New("illegal plant state transition")
	ErrSnapshotStale    = errors.New("snapshot revision stale")
	ErrWindowOpen       = errors.New("timing window still open")
	ErrPreheatIncomplete  = errors.New("reheatfn preheat incomplete")
	ErrCoordinationLock = errors.New("coordination lock held")
	ErrSlabskidLevelLow     = errors.New("slabskid level below low limit")
	ErrReheatLoss        = errors.New("reheatfn reheat lost")
	ErrDescalcLimit    = errors.New("descalc valve at limit")
)
