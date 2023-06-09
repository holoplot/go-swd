package io

import (
	"testing"
)

func TestTransaction_StartByte(t *testing.T) {
	type fields struct {
		AccessPort bool
		Write      bool
		Address    uint8
		Data       uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "1",
			fields: fields{
				AccessPort: false,
				Write:      false,
				Address:    0,
			},
			want: 0x81,
		},
		{
			name: "2",
			fields: fields{
				AccessPort: true,
				Write:      false,
				Address:    0,
			},
			want: 0xa3,
		},
		{
			name: "3",
			fields: fields{
				AccessPort: true,
				Write:      true,
				Address:    3,
			},
			want: 0x9f,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				AccessPort: tt.fields.AccessPort,
				Write:      tt.fields.Write,
				Address:    tt.fields.Address,
				Data:       tt.fields.Data,
			}
			if got := tx.StartByte(); got != tt.want {
				t.Errorf("Transaction.StartByteParity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTransaction_DataParity(t *testing.T) {
	type fields struct {
		AccessPort bool
		Write      bool
		Address    uint8
		Data       uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   uint8
	}{
		{
			name: "1",
			fields: fields{
				Data: 0x00000000,
			},
			want: 0,
		},
		{
			name: "2",
			fields: fields{
				Data: 0x00010000,
			},
			want: 1,
		},
		{
			name: "3",
			fields: fields{
				Data: 0xffffffff,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				AccessPort: tt.fields.AccessPort,
				Write:      tt.fields.Write,
				Address:    tt.fields.Address,
				Data:       tt.fields.Data,
			}
			if got := tx.DataParity(); got != tt.want {
				t.Errorf("Transaction.DataParity() = %v, want %v", got, tt.want)
			}
		})
	}
}
