package interlock

import (
	"fmt"

	"github.com/lacsar712/slabheat/internal/model"
)

type PermissiveSet struct {
	gasfuelOK       bool
	ignitionOK   bool
	slabskidOK       bool
	pressureOK   bool
	burnerOK bool
}

func NewPermissiveSet() *PermissiveSet { return &PermissiveSet{} }

func (p *PermissiveSet) SetGasfuel(ok bool)       { p.gasfuelOK = ok }
func (p *PermissiveSet) SetIgnition(ok bool)   { p.ignitionOK = ok }
func (p *PermissiveSet) SetSlabskid(ok bool)       { p.slabskidOK = ok }
func (p *PermissiveSet) SetPressure(ok bool)   { p.pressureOK = ok }
func (p *PermissiveSet) SetBurner(ok bool) { p.burnerOK = ok }

func (p *PermissiveSet) GasfuelOK() bool       { return p.gasfuelOK }
func (p *PermissiveSet) IgnitionOK() bool   { return p.ignitionOK }
func (p *PermissiveSet) SlabskidOK() bool       { return p.slabskidOK }
func (p *PermissiveSet) PressureOK() bool   { return p.pressureOK }
func (p *PermissiveSet) BurnerOK() bool { return p.burnerOK }

func (p *PermissiveSet) AllFiring() bool {
	return p.gasfuelOK && p.ignitionOK && p.slabskidOK && p.pressureOK && p.burnerOK
}

func (p *PermissiveSet) CheckIgnition() error {
	if !p.gasfuelOK {
		return fmt.Errorf("%w", model.ErrGasfuelPermissive)
	}
	if !p.ignitionOK {
		return fmt.Errorf("%w", model.ErrIgnitionBlocked)
	}
	return nil
}

func CheckReheatLoss(reading model.BurnerReading) error {
	if reading.BurnerPhase == model.BurnerStable && reading.ReheatfnTempF < 600 {
		return fmt.Errorf("%w", model.ErrReheatLoss)
	}
	return nil
}

func (p *PermissiveSet) CheckFiring() error {
	if err := p.CheckIgnition(); err != nil {
		return err
	}
	if !p.slabskidOK {
		return fmt.Errorf("%w", model.ErrSlabskidLevelTrip)
	}
	if !p.pressureOK {
		return fmt.Errorf("%w", model.ErrPressureTrip)
	}
	if !p.burnerOK {
		return fmt.Errorf("%w", model.ErrBurnerTrip)
	}
	return nil
}

type CoordinationLock struct {
	holder string
	held   bool
}

func NewCoordinationLock() *CoordinationLock { return &CoordinationLock{} }

func (c *CoordinationLock) Acquire(holder string) error {
	if c.held {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	c.holder = holder
	c.held = true
	return nil
}

func (c *CoordinationLock) Release(holder string) {
	if c.held && c.holder == holder {
		c.held = false
		c.holder = ""
	}
}

func (c *CoordinationLock) Require(holder string) error {
	if !c.held || c.holder != holder {
		return fmt.Errorf("%w", model.ErrCoordinationLock)
	}
	return nil
}

func (c *CoordinationLock) Held() bool { return c.held }
