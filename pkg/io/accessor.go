package io

type Accessor interface {
	Reset() error
	Tx(*Transaction) error
	Close()
}
