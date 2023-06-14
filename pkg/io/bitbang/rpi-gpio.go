//go:build linux
// +build linux

package bitbang

import (
	"github.com/stianeikeland/go-rpio"
)

type rpiGPIOhw struct {
	clk rpio.Pin
	io  rpio.Pin
}

func (l *rpiGPIOhw) SetClock(v int) error {
	if v == 0 {
		l.clk.Low()
	} else {
		l.clk.High()
	}

	return nil
}

func (l *rpiGPIOhw) SetData(v int) error {
	if v == 0 {
		l.io.Low()
	} else {
		l.io.High()
	}

	return nil
}

func (l *rpiGPIOhw) GetData() (int, error) {
	return int(l.io.Read()), nil
}

func (l *rpiGPIOhw) SetDataDirectionInput() error {
	l.io.Input()

	return nil
}

func (l *rpiGPIOhw) SetDataDirectionOutput() error {
	l.io.Output()

	return nil
}

func (l *rpiGPIOhw) Close() {
	l.io.Input()
	l.clk.Input()
}

func NewRPI(chip string, ioGPIO, clkGPIO, frequency int) (*BitBang, error) {
	hw := &rpiGPIOhw{
		io:  rpio.Pin(ioGPIO),
		clk: rpio.Pin(clkGPIO),
	}

	return New(hw, frequency), nil
}
