//go:build linux
// +build linux

package bitbang

import (
	"github.com/warthog618/gpiod"
)

type linuxGPIOhw struct {
	io  *gpiod.Line
	clk *gpiod.Line
}

func (l *linuxGPIOhw) SetClock(v int) error {
	return l.clk.SetValue(v)
}

func (l *linuxGPIOhw) SetData(v int) error {
	return l.io.SetValue(v)
}

func (l *linuxGPIOhw) GetData() (int, error) {
	return l.io.Value()
}

func (l *linuxGPIOhw) SetDataDirectionInput() error {
	return l.io.Reconfigure(gpiod.AsInput)
}

func (l *linuxGPIOhw) SetDataDirectionOutput() error {
	return l.io.Reconfigure(gpiod.AsOutput(0))
}

func (l *linuxGPIOhw) Close() {
	l.io.Close()
	l.clk.Close()
}

func NewLinuxGPIO(chip string, ioGPIO, clkGPIO, frequency int) (*BitBang, error) {
	c, err := gpiod.NewChip(chip)
	if err != nil {
		return nil, err
	}

	ioLine, err := c.RequestLine(ioGPIO, gpiod.AsInput)
	if err != nil {
		return nil, err
	}

	clkLine, err := c.RequestLine(clkGPIO, gpiod.AsOutput(0))
	if err != nil {
		return nil, err
	}

	hw := &linuxGPIOhw{
		io:  ioLine,
		clk: clkLine,
	}

	return New(hw, frequency), nil
}
