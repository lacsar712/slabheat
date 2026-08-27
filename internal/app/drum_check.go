package app

import (
	"fmt"

	"github.com/lacsar712/slabheat/internal/model"
)

func (a *App) CheckSlabskidLevel(snap model.PlantSnapshot) error {
	if snap.Slabskid.LevelPercent < model.MinSlabskidLevelPercent {
		return fmt.Errorf("level low: %.1f", snap.Slabskid.LevelPercent)
	}
	return nil
}
