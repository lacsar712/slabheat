package model

import (
	"errors"
	"fmt"
)

func LevelChain(pct float64) error {
	return fmt.Errorf("level %.1f: %w", pct, ErrSlabskidLevelLow)
}

func PermissiveChain(kind string) error {
	return fmt.Errorf("permissive %s: %w", kind, ErrGasfuelPermissive)
}

func LossChain(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("thermal loss: %w", err)
}

func IsLevel(err error) bool { return errors.Is(err, ErrSlabskidLevelLow) }
func IsPermissive(err error) bool { return errors.Is(err, ErrGasfuelPermissive) }
func IsLoss(err error) bool { return errors.Is(err, ErrReheatLoss) }
