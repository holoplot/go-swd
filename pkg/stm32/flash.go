package stm32

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/holoplot/go-swd/pkg/swd"
	"github.com/rs/zerolog/log"
)

const (
	flashBaseAddr uint32 = 0x08000000

	regBase    uint32 = 0x40022000
	regACR     uint32 = regBase + 0x00
	regKEYR    uint32 = regBase + 0x08
	regOPTKEYR uint32 = regBase + 0x0c
	regSR      uint32 = regBase + 0x10
	regCR      uint32 = regBase + 0x14
	regECCR    uint32 = regBase + 0x18
	regOPTR    uint32 = regBase + 0x20

	flashKey1 = 0x45670123
	flashKey2 = 0xcdef89ab
)

type acrRegister uint32

const (
	acrRegisterLatencyMask       acrRegister = 0x7
	acrRegisterCPUPrefetchEnable acrRegister = 1 << 8
	acrRegisterInstructionCache  acrRegister = 1 << 9
	acrRegisterDataCache         acrRegister = 1 << 11
	acrEmpty                     acrRegister = 1 << 16
	acrDebugEnable               acrRegister = 1 << 18
)

type statusRegister uint32

const (
	statusRegisterEndOfOperation            statusRegister = 1 << 0
	statusRegisterOperationError            statusRegister = 1 << 1
	statusRegisterProgrammingError          statusRegister = 1 << 3
	statusRegisterWriteProtectionError      statusRegister = 1 << 4
	statusRegisterProgrammingAlignmentError statusRegister = 1 << 5
	statusRegisterSizeError                 statusRegister = 1 << 6
	statusRegisterProgrammingSequenceError  statusRegister = 1 << 7
	statusRegisterStatusMissError           statusRegister = 1 << 8
	statusRegisterFastProgrammingError      statusRegister = 1 << 9
	statusRegisterBusy1                     statusRegister = 1 << 16
	statusRegisterBusy2                     statusRegister = 1 << 17
)

type controlRegister uint32

const (
	controlRegisterPg                    controlRegister = 1 << 0
	controlRegisterPer                   controlRegister = 1 << 1
	controlRegisterMer1                  controlRegister = 1 << 2
	controlRegisterMer2                  controlRegister = 1 << 15
	controlRegisterStart                 controlRegister = 1 << 16
	controlRegisterOptionStart           controlRegister = 1 << 17
	controlRegisterOptionFastProgramming controlRegister = 1 << 18
	controlRegisterEndOfOperation        controlRegister = 1 << 24
	controlRegisterOptLock               controlRegister = 1 << 30
	controlRegisterLock                  controlRegister = 1 << 31
)

type Flash struct {
	swd        *swd.SWD
	isWritable bool
	pllEnabled bool
}

func (f *Flash) Read(addr, size uint32, writer io.Writer) error {
	for i := uint32(0); i < size; i += 4 {
		data, err := f.swd.ReadRegister(flashBaseAddr + addr + i)
		if err != nil {
			return err
		}

		if err := binary.Write(writer, binary.LittleEndian, data); err != nil {
			return err
		}
	}

	return nil
}

func (f *Flash) makeWriteable() error {
	if !f.pllEnabled {
		log.Info().Msgf("Enable PLL!!")

		if err := f.swd.WriteRegister(0x40021008, 0x00000002); err != nil {
			return err
		}

		if err := f.swd.WriteRegister(0x4002100c, 0xF60A1812); err != nil {
			return err
		}

		if err := f.swd.WriteRegister(0x40021000, 0x1000500); err != nil {
			return err
		}

		// for {
		// 	val, err := f.swd.ReadRegister(0x40021000)
		// 	if err != nil {
		// 		return err
		// 	}

		// 	log.Info().Msgf("RCC CR: 0x%08x", val)

		// 	if val&(1<<25) != 0 {
		// 		break
		// 	}
		// }

		f.pllEnabled = true
	}

	if !f.isWritable {
		if err := f.swd.WriteRegister(regKEYR, flashKey1); err != nil {
			return err
		}

		if err := f.swd.WriteRegister(regKEYR, flashKey2); err != nil {
			return err
		}

		f.isWritable = true
	}

	return nil
}

func (f *Flash) clearErrors() error {
	clr := statusRegisterOperationError |
		statusRegisterProgrammingError |
		statusRegisterWriteProtectionError |
		statusRegisterProgrammingAlignmentError |
		statusRegisterSizeError |
		statusRegisterProgrammingSequenceError |
		statusRegisterStatusMissError |
		statusRegisterFastProgrammingError

	return f.swd.WriteRegister(regSR, uint32(clr))
}

func (f *Flash) waitForCompletion() error {
	var sr statusRegister

	defer func() {
		_ = f.swd.WriteRegister(regSR, uint32(sr))
	}()

	for {
		val, err := f.swd.ReadRegister(regSR)
		if err != nil {
			return err
		}

		sr = statusRegister(val)

		log.Info().Msgf("SR: 0x%08x", sr)

		if sr&statusRegisterProgrammingAlignmentError != 0 {
			return fmt.Errorf("programming alignment error")
		}

		if sr&statusRegisterProgrammingError != 0 {
			return fmt.Errorf("programming error")
		}

		if sr&statusRegisterOperationError != 0 {
			return fmt.Errorf("operation error")
		}

		if sr&(statusRegisterBusy1|statusRegisterEndOfOperation) == statusRegisterEndOfOperation {
			panic("end of operation")

			return nil
		}

		time.Sleep(time.Millisecond)
	}
}

func (f *Flash) Write(addr uint32, reader io.Reader) error {
	if err := f.makeWriteable(); err != nil {
		return fmt.Errorf("make writable: %w", err)
	}

	if err := f.swd.UpdateCSW(swd.CSWAutoIncrementOff|swd.CSWSize32bit,
		swd.CSWAutoIncrementMask|swd.CSWSizeMask); err != nil {
		return err
	}

	for {
		if busy, err := f.busy(); err != nil {
			return err
		} else if !busy {
			break
		}
	}

	if err := f.clearErrors(); err != nil {
		return err
	}

	if err := f.swd.WriteRegister(regCR, uint32(controlRegisterPg|controlRegisterEndOfOperation)); err != nil {
		return err
	}

	defer func() {
		_ = f.swd.WriteRegister(regCR, 0)
	}()

	for {
		var data1, data2 uint32

		if err := binary.Read(reader, binary.LittleEndian, &data1); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}

		// Write the second word as 0 in case of errors
		_ = binary.Read(reader, binary.LittleEndian, &data2)

		if err := f.swd.WriteTAR(flashBaseAddr + addr); err != nil {
			return err
		}

		if err := f.swd.WriteDRW(data1); err != nil {
			return err
		}

		addr += 4

		if err := f.swd.WriteTAR(flashBaseAddr + addr); err != nil {
			return err
		}

		if err := f.swd.WriteDRW(uint32(data2)); err != nil {
			return err
		}

		addr += 4

		if err := f.waitForCompletion(); err != nil {
			return err
		}
	}

	log.Info().Msgf("Flash write done")

	// Clear EMPTY bit
	// if err := f.swd.UpdateRegisterBits(regACR, uint32(acrEmpty|acrRegisterLatencyMask), 2); err != nil {
	// 	return err
	// }
	if err := f.swd.WriteRegister(regACR, 2); err != nil {
		return err
	}

	return nil
}

var ErrTimeout = errors.New("timeout")

func (f *Flash) EraseAll(timeout time.Duration) error {
	if err := f.makeWriteable(); err != nil {
		return fmt.Errorf("failed to make writable: %w", err)
	}

	for {
		if busy, err := f.busy(); err != nil {
			return err
		} else if !busy {
			break
		}
	}

	if err := f.clearErrors(); err != nil {
		return err
	}

	if err := f.swd.WriteRegister(regCR, uint32(controlRegisterMer1|controlRegisterStart)); err != nil {
		return err
	}

	start := time.Now()

	for time.Since(start) < timeout {
		if busy, err := f.busy(); err != nil {
			return err
		} else if !busy {
			return nil
		}

		time.Sleep(time.Millisecond * 100)
	}

	return ErrTimeout
}

func (f *Flash) Initialize() error {
	if busy, err := f.busy(); err != nil {
		return err
	} else if busy {
		return fmt.Errorf("flash is busy")
	}

	// if err := f.swd.UpdateRegisterBits(regACR, uint32(acrRegisterLatencyMask), 2); err != nil {
	// 	return err
	// }
	if err := f.swd.WriteRegister(regACR, 2); err != nil {
		return err
	}

	cr, err := f.swd.ReadRegister(regCR)
	if err != nil {
		return err
	}

	f.isWritable = (cr & uint32(controlRegisterLock)) == 0

	return nil
}

func (f *Flash) busy() (bool, error) {
	sr, err := f.swd.ReadRegister(regSR)
	if err != nil {
		return false, err
	}

	return (statusRegister(sr) & statusRegisterBusy1) != 0, nil
}

func newFlash(swd *swd.SWD) *Flash {
	return &Flash{
		swd: swd,
	}
}
