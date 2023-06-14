//go:build linux
// +build linux

package bitbang

import (
	"fmt"

	"github.com/warthog618/gpiod"
)

type linuxGPIOhw struct {
	io  *gpiod.Line
	clk *gpiod.Line
}

func (l *linuxGPIOhw) SetClock(v int) error {
	if err := l.clk.SetValue(v); err != nil {
		return fmt.Errorf("error setting clock line: %w", err)
	}

	return nil
}

func (l *linuxGPIOhw) SetData(v int) error {
	if err := l.io.SetValue(v); err != nil {
		return fmt.Errorf("error setting data line: %w", err)
	}

	return nil
}

func (l *linuxGPIOhw) GetData() (int, error) {
	if v, err := l.io.Value(); err != nil {
		return 0, fmt.Errorf("error reading data line: %w", err)
	} else {
		return v, nil
	}
}

func (l *linuxGPIOhw) SetDataDirectionInput() error {
	if err := l.io.Reconfigure(gpiod.AsInput); err != nil {
		return fmt.Errorf("error setting data line to output direction: %w", err)
	}

	return nil
}

func (l *linuxGPIOhw) SetDataDirectionOutput() error {
	if err := l.io.Reconfigure(gpiod.AsOutput(0)); err != nil {
		return fmt.Errorf("error setting data line to input direction: %w", err)
	}

	return nil
}

func (l *linuxGPIOhw) Close() {
	_ = l.io.Reconfigure(gpiod.AsInput)
	_ = l.clk.Reconfigure(gpiod.AsInput)

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
