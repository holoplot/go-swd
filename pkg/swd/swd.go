package swd

import (
	"github.com/holoplot/go-swd/pkg/io"
)

const (
	regIDCODE = 0x0
)

type SWD struct {
	accessor io.Accessor
}

func (s *SWD) IDCode() (uint32, error) {
	tx := &io.Transaction{
		AccessPort: true,
		Write:      false,
		Address:    regIDCODE,
	}

	if err := s.accessor.Tx(tx); err != nil {
		return 0, err
	}

	return tx.Data, nil
}

func New(accessor io.Accessor) *SWD {
	return &SWD{
		accessor: accessor,
	}
}
