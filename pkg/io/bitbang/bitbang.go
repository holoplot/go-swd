package bitbang

import (
	"time"

	"github.com/holoplot/go-swd/pkg/io"
)

type BitBanger interface {
	SetClock(int) error
	SetData(int) error
	GetData() (int, error)
	SetDataDirectionInput() error
	SetDataDirectionOutput() error
	Close()
}

type BitBang struct {
	hw             BitBanger
	clockHalfCycle time.Duration
}

func (bb *BitBang) read(n int) (uint32, error) {
	var v uint32

	for i := 0; i < n; i++ {
		if err := bb.hw.SetClock(0); err != nil {
			return 0, err
		}

		x, err := bb.hw.GetData()
		if err != nil {
			return 0, err
		}

		v |= uint32(x) << i

		time.Sleep(bb.clockHalfCycle)

		if err := bb.hw.SetClock(1); err != nil {
			return 0, err
		}

		time.Sleep(bb.clockHalfCycle)
	}

	return v, nil
}

func (bb *BitBang) write(v uint32, n int) error {
	for i := 0; i < n; i++ {
		if err := bb.hw.SetClock(0); err != nil {
			return err
		}

		if err := bb.hw.SetData(int(v & 1)); err != nil {
			return err
		}

		v >>= 1

		time.Sleep(bb.clockHalfCycle)

		if err := bb.hw.SetClock(1); err != nil {
			return err
		}

		time.Sleep(bb.clockHalfCycle)
	}

	return nil
}

func (bb *BitBang) dummyClockCycle() error {
	if err := bb.hw.SetClock(0); err != nil {
		return err
	}

	time.Sleep(bb.clockHalfCycle)

	if err := bb.hw.SetClock(1); err != nil {
		return err
	}

	time.Sleep(bb.clockHalfCycle)

	return nil
}

func (bb *BitBang) LineReset() error {
	if err := bb.hw.SetDataDirectionOutput(); err != nil {
		return err
	}

	if err := bb.hw.SetData(1); err != nil {
		return err
	}

	for i := 0; i < 54; i++ {
		if err := bb.dummyClockCycle(); err != nil {
			return err
		}
	}

	if err := bb.write(0, 1); err != nil {
		return err
	}

	return nil
}

func (bb *BitBang) Tx(tx *io.Transaction) error {
	if err := bb.hw.SetDataDirectionOutput(); err != nil {
		return err
	}

	if err := bb.write(uint32(tx.RequestByte()), 8); err != nil {
		return err
	}

	if err := bb.hw.SetDataDirectionInput(); err != nil {
		return err
	}

	if err := bb.dummyClockCycle(); err != nil {
		return err
	}

	ack, err := bb.read(3)
	if err != nil {
		return err
	}

	tx.Ack = io.Ack(ack)

	if tx.Ack != io.AckOk && tx.Ack != io.AckWait {
		if err := bb.dummyClockCycle(); err != nil {
			return err
		}

		return io.ErrBadAck
	}

	if tx.Direction == io.DirectionWrite {
		// Turnaround
		if err := bb.dummyClockCycle(); err != nil {
			return err
		}

		if err := bb.hw.SetDataDirectionOutput(); err != nil {
			return err
		}

		if err := bb.write(tx.Data, 32); err != nil {
			return err
		}

		if err := bb.write(uint32(tx.DataParity().Bit()), 1); err != nil {
			return err
		}
	} else {
		data, err := bb.read(32)
		if err != nil {
			return err
		}

		tx.Data = data

		parity, err := bb.read(1)
		if err != nil {
			return err
		}

		if parity != tx.DataParity().Bit() {
			return io.ErrBadParity
		}
	}

	if err := bb.hw.SetDataDirectionOutput(); err != nil {
		return err
	}

	// dummy write to keep the interface running for at least 8 cycles
	if err := bb.write(0, 8); err != nil {
		return err
	}

	return nil
}

func (bb *BitBang) Close() {
	bb.hw.Close()
}

func New(hw BitBanger, frequency int) *BitBang {
	return &BitBang{
		hw:             hw,
		clockHalfCycle: time.Second / time.Duration(frequency) / 2,
	}
}
