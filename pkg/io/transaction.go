package io

import (
	"errors"
	"fmt"
)

var (
	ErrBadParity = errors.New("bad parity")
	ErrBadAck    = errors.New("bad ack")
)

type RequestByte byte

const (
	startBit  RequestByte = 1 << 0
	apBit     RequestByte = 1 << 1
	rdBit     RequestByte = 1 << 2
	addrShift RequestByte = 3
	addrMask  RequestByte = 3 << addrShift
	parityBit RequestByte = 1 << 5
	stopBit   RequestByte = 0 << 6
	parkBit   RequestByte = 1 << 7
)

type PortType bool

const (
	DebugPort  PortType = false
	AccessPort PortType = true
)

func (p PortType) requestByteShifted() RequestByte {
	if p == AccessPort {
		return apBit
	}

	return 0
}

func (p PortType) String() string {
	if p == AccessPort {
		return "AP"
	}

	return "DP"
}

type Direction bool

const (
	DirectionWrite Direction = false
	DirectionRead  Direction = true
)

func (d Direction) requestByteShifted() RequestByte {
	if d == DirectionRead {
		return rdBit
	}

	return 0
}

func (d Direction) String() string {
	if d == DirectionRead {
		return "read"
	}

	return "write"
}

type Address uint8

func (a Address) requestByteShifted() RequestByte {
	return (RequestByte(a>>2) << addrShift) & addrMask
}

func (a Address) String() string {
	return fmt.Sprintf("0x%1x", uint8(a))
}

type Ack uint8

const (
	AckOk    Ack = 0b001
	AckWait  Ack = 0b010
	AckFault Ack = 0b100
)

func (ack Ack) String() string {
	switch ack {
	case AckOk:
		return "ok"
	case AckWait:
		return "wait"
	case AckFault:
		return "fault"
	default:
		return fmt.Sprintf("unknown:%d", ack)
	}
}

type Transaction struct {
	PortType  PortType
	Direction Direction
	Address   Address
	Data      uint32
	Ack       Ack
}

type Parity bool

const (
	ParityEven Parity = false
	ParityOdd  Parity = true
)

func (p Parity) Bit() uint32 {
	if p == ParityOdd {
		return 1
	}

	return 0
}

func ParityFromBit(b uint32) Parity {
	if b == 1 {
		return ParityOdd
	}

	return ParityEven
}

func (p Parity) requestByteShifted() RequestByte {
	if p == ParityOdd {
		return parityBit
	}

	return 0
}

func (tx *Transaction) RequestByte() RequestByte {
	v := startBit | stopBit | parkBit |
		tx.Direction.requestByteShifted() |
		tx.PortType.requestByteShifted() |
		tx.Address.requestByteShifted()

	parity := []Parity{
		ParityEven, // 0b0000
		ParityOdd,  // 0b0001
		ParityOdd,  // 0b0010
		ParityEven, // 0b0011

		ParityOdd,  // 0b0100
		ParityEven, // 0b0101
		ParityEven, // 0b0110
		ParityOdd,  // 0b0111

		ParityOdd,  // 0b1000
		ParityEven, // 0b1001
		ParityEven, // 0b1010
		ParityOdd,  // 0b1011

		ParityEven, // 0b1100
		ParityOdd,  // 0b1101
		ParityOdd,  // 0b1110
		ParityEven, // 0b1111
	}

	// shift out start bit; we're only interested in bits 1-4
	v |= parity[(v>>1)&0xf].requestByteShifted()

	return v
}

func (tx *Transaction) DataParity() Parity {
	i := tx.Data
	i = i - ((i >> 1) & 0x55555555)
	i = (i & 0x33333333) + ((i >> 2) & 0x33333333)
	i = (((i + (i >> 4)) & 0x0f0f0f0f) * 0x01010101) >> 24

	return Parity(i&1 == 1)
}
