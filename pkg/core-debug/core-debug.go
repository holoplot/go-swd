package coredebug

import (
	"errors"
	"fmt"
	"time"

	"github.com/holoplot/go-swd/pkg/swd"
	"github.com/rs/zerolog/log"
)

const (
	retries = 10
)

type CoreDebug struct {
	swd *swd.SWD
}

func (cd *CoreDebug) ReadDHCSR() (DHCSR, error) {
	reg, err := cd.swd.ReadRegister(regDHCSR)
	if err != nil {
		return 0, err
	}

	return DHCSR(reg), nil
}

func (cd *CoreDebug) WriteDHCSR(dhcsr DHCSR) error {
	return cd.swd.WriteRegister(regDHCSR, uint32(dhcsr))
}

func (cd *CoreDebug) ReadDCRSR() (DCRSR, error) {
	reg, err := cd.swd.ReadRegister(regDCRSR)
	if err != nil {
		return 0, err
	}

	return DCRSR(reg), nil
}

func (cd *CoreDebug) WriteDCRSR(dcrsr DCRSR) error {
	return cd.swd.WriteRegister(regDCRSR, uint32(dcrsr))
}

func (cd *CoreDebug) ReadDCRDR() (DCRDR, error) {
	reg, err := cd.swd.ReadRegister(regDCRDR)
	if err != nil {
		return 0, err
	}

	return DCRDR(reg), nil
}

func (cd *CoreDebug) ReadDEMCR() (DEMCR, error) {
	reg, err := cd.swd.ReadRegister(regDEMCR)
	if err != nil {
		return 0, err
	}

	return DEMCR(reg), nil
}

func (cd *CoreDebug) WriteDEMCR(demcr DEMCR) error {
	return cd.swd.WriteRegister(regDEMCR, uint32(demcr))
}

func (cd *CoreDebug) WriteDCRDR(dcrdr DCRDR) error {
	return cd.swd.WriteRegister(regDCRDR, uint32(dcrdr))
}

func (cd *CoreDebug) ResetRegisters() error {
	if err := cd.WriteDHCSR(DHCSRDebugKey); err != nil {
		return err
	}

	if err := cd.WriteDEMCR(0); err != nil {
		return err
	}

	return nil
}

var ErrTimeout = errors.New("timeout")

func (cd *CoreDebug) Halt() error {
	if err := cd.WriteDHCSR(DHCSRDebugKey | DHCSRCDebugEn | DHCSRCHalt); err != nil {
		return err
	}

	if err := cd.WriteDHCSR(DHCSRDebugKey | DHCSRCDebugEn | DHCSRCHalt | DHCSRCMaskInts); err != nil {
		return err
	}

	for n := 0; n < retries; n++ {
		dhcsr, err := cd.ReadDHCSR()
		if err != nil {
			return fmt.Errorf("error reading DHCSR: %w", err)
		}

		log.Info().Msgf("DHCSR: 0x%08x", dhcsr)

		if dhcsr&DHCSRSHalt != 0 {
			return nil
		}

		time.Sleep(time.Millisecond * 100)
	}

	return ErrTimeout
}

func (cd *CoreDebug) Continue() error {
	if err := cd.WriteDHCSR(DHCSRDebugKey | DHCSRCDebugEn); err != nil {
		return err
	}

	for n := 0; n < retries; n++ {
		dhcsr, err := cd.ReadDHCSR()
		if err != nil {
			return fmt.Errorf("error reading DHCSR: %w", err)
		}

		if dhcsr&DHCSRSHalt == 0 {
			return nil
		}

		time.Sleep(time.Millisecond)
	}

	return ErrTimeout
}

func (cd *CoreDebug) RunAfterReset() error {
	if err := cd.WriteDHCSR(DHCSRDebugKey); err != nil {
		return err
	}

	return nil

	for n := 0; n < retries; n++ {
		dhcsr, err := cd.ReadDHCSR()
		if err != nil {
			return fmt.Errorf("error reading DHCSR: %w", err)
		}

		if dhcsr&DHCSRSHalt == 0 {
			return nil
		}

		time.Sleep(time.Millisecond)
	}

	return ErrTimeout
}

func New(swd *swd.SWD) *CoreDebug {
	return &CoreDebug{
		swd: swd,
	}
}
