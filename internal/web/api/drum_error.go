package api

import (
	"errors"

	"github.com/lacsar712/slabheat/internal/model"
)

func classifySlabskidError(err error) (string, bool) {
	if errors.Is(err, model.ErrSlabskidLevelLow) {
		return "slabskid_level_low", true
	}
	return "", false
}
