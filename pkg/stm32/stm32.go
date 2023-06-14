package stm32

import (
	cd "github.com/holoplot/go-swd/pkg/core-debug"
	"github.com/holoplot/go-swd/pkg/swd"
	scb "github.com/holoplot/go-swd/pkg/system-control-block"
)

type STM32 struct {
	swd       *swd.SWD
	flash     *Flash
	coreDebug *cd.CoreDebug
	scb       *scb.SystemControlBlock
}

func (stm *STM32) Flash() *Flash {
	return stm.flash
}

func (stm *STM32) Reset() error {
	if err := stm.coreDebug.ResetRegisters(); err != nil {
		return err
	}

	return stm.scb.ResetSystem()
}

func (stm *STM32) Halt() error {
	return stm.coreDebug.Halt()
}

func (stm *STM32) RunAfterReset() error {
	return stm.coreDebug.RunAfterReset()
}

func New(swd *swd.SWD) *STM32 {
	stm32 := &STM32{
		swd:       swd,
		flash:     newFlash(swd),
		coreDebug: cd.New(swd),
		scb:       scb.New(swd),
	}

	return stm32
}
