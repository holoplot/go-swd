package io

import "errors"

var (
	ErrBadParity = errors.New("bad parity")
	ErrBadAck    = errors.New("bad ack")
)

const (
	AckOk = 1
)

type Transaction struct {
	// DP/AP
	AccessPort bool
	// Read/Write
	Write bool
	// A2/A3
	Address uint8

	Data uint32
}

func (tx *Transaction) StartByte() uint8 {
	boolToUint8 := func(b bool) uint8 {
		if b {
			return 1
		}

		return 0
	}

	parity := boolToUint8(tx.AccessPort) ^
		boolToUint8(tx.Write) ^
		(tx.Address & 1) ^
		((tx.Address >> 1) & 1)

	return (1 << 0) | // start bit
		(boolToUint8(tx.AccessPort) << 1) |
		(boolToUint8(tx.Write) << 2) |
		(tx.Address << 3) |
		(parity << 5) |
		(0 << 6) | // stop bit
		(1 << 7) // park bit
}

func (tx *Transaction) DataParity() uint8 {
	i := tx.Data
	i = i - ((i >> 1) & 0x55555555)
	i = (i & 0x33333333) + ((i >> 2) & 0x33333333)
	i = (((i + (i >> 4)) & 0x0f0f0f0f) * 0x01010101) >> 24

	return uint8(i & 1)
}
