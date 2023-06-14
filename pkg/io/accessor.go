package io

type Accessor interface {
	LineReset() error
	// Read(RequestByte) (uint32, Ack, Parity, error)
	// Write(RequestByte, uint32, Parity) (Ack, error)
	Tx(*Transaction) error
	Close()
}
