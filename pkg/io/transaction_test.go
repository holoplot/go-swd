package io

import (
	"testing"
)

func TestRequestByte(t *testing.T) {
	tests := []struct {
		name string
		tx   *Transaction
		want RequestByte
	}{
		{
			name: "1",
			tx: &Transaction{
				PortType:  DebugPort,
				Direction: DirectionWrite,
				Address:   0x00,
			},
			want: 0x81,
		},
		{
			name: "2",
			tx: &Transaction{
				PortType:  DebugPort,
				Direction: DirectionWrite,
				Address:   0x04,
			},
			want: 0xa9,
		},
		{
			name: "3",
			tx: &Transaction{
				PortType:  DebugPort,
				Direction: DirectionRead,
				Address:   0x0c,
			},
			want: 0xbd,
		},
		{
			name: "4",
			tx: &Transaction{
				PortType:  AccessPort,
				Direction: DirectionRead,
				Address:   0x0c,
			},
			want: 0x9f,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tx.RequestByte(); got != tt.want {
				t.Errorf("Transaction.RequestByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_DataParity(t *testing.T) {
	type fields struct {
		Data uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   Parity
	}{
		{
			name: "1",
			fields: fields{
				Data: 0x00000000,
			},
			want: ParityEven,
		},
		{
			name: "2",
			fields: fields{
				Data: 0x00010000,
			},
			want: ParityOdd,
		},
		{
			name: "3",
			fields: fields{
				Data: 0xffffffff,
			},
			want: ParityEven,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				Data: tt.fields.Data,
			}
			if got := tx.DataParity(); got != tt.want {
				t.Errorf("Transaction.DataParity() = %v, want %v", got, tt.want)
			}
		})
	}
}
