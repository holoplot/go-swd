package debug

import "github.com/holoplot/go-swd/pkg/io"

type Debugger interface {
	Tx(string, io.Transaction, error)
}

type NopDebugger struct{}

func (d *NopDebugger) Tx(string, io.Transaction, error) {}
