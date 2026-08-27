package app

import (
	"fmt"

	"github.com/lacsar712/slabheat/internal/model"
)

func (a *App) CheckSlabskidLevel(snap model.PlantSnapshot) error {
	if snap.Slabskid.LevelPercent < model.MinSlabskidLevelPercent {
		return model.LevelChain(snap.Slabskid.LevelPercent)
	}
	return nil
}
