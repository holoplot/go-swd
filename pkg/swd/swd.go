package swd

import (
	"errors"
	"fmt"
	"time"

	"github.com/holoplot/go-swd/pkg/debug"
	"github.com/holoplot/go-swd/pkg/io"
)

var (
	ErrTimeout = errors.New("timeout")
)

type SWD struct {
	accessor io.Accessor
	debugger debug.Debugger

	currentSelect uint32
}

func (s *SWD) writeTx(name string, portType io.PortType, addr io.Address, data uint32) error {
	for {
		tx := &io.Transaction{
			PortType:  portType,
			Direction: io.DirectionWrite,
			Address:   addr,
			Data:      data,
		}

		err := s.accessor.Tx(tx)

		s.debugger.Tx(name, *tx, err)

		if err != nil {
			return err
		}

		if tx.Ack == io.AckOk {
			return nil
		}

		// Repeat on AckWait
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *SWD) readTx(name string, portType io.PortType, addr io.Address) (uint32, error) {
	for {
		tx := &io.Transaction{
			PortType:  portType,
			Direction: io.DirectionRead,
			Address:   addr,
		}

		err := s.accessor.Tx(tx)

		s.debugger.Tx(name, *tx, err)

		if err != nil {
			return 0, err
		}

		if tx.Ack == io.AckOk {
			return tx.Data, nil
		}

		// Repeat on AckWait
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *SWD) Abort(flags AbortFlags) error {
	return s.writeTx("ABORT", io.DebugPort, regAbort, uint32(flags))
}

func (s *SWD) Select(accessPort uint32, bank uint8, low uint8) error {
	v := (accessPort << 24) | (uint32(bank) << 4) | uint32(low)

	if s.currentSelect == v {
		return nil
	}

	if err := s.writeTx("SELECT", io.DebugPort, regSelect, v); err != nil {
		return fmt.Errorf("select failed: %w", err)
	}

	s.currentSelect = v

	return nil
}

func (s *SWD) WriteMemAP(name string, addr io.Address, data uint32) error {
	if err := s.Select(0, uint8(addr>>4), 0); err != nil {
		return err
	}

	name = fmt.Sprintf("MEMAP:%s", name)

	return s.writeTx(name, io.AccessPort, io.Address(addr&0xf), data)
}

func (s *SWD) ReadRdBuff() (uint32, error) {
	return s.readTx("RDBUFF", io.DebugPort, regReadBuffer)
}

func (s *SWD) ReadMemAP(name string, addr io.Address) (uint32, error) {
	if err := s.Select(0, uint8(addr>>4), 0); err != nil {
		return 0, fmt.Errorf("select: %w", err)
	}

	name = fmt.Sprintf("MEMAP:%s", name)

	if _, err := s.readTx(name, io.AccessPort, io.Address(uint8(addr)&0xf)); err != nil {
		return 0, fmt.Errorf("read tx: %w, addr %02x", err, uint8(addr)&0xf)
	}

	v, err := s.ReadRdBuff()
	if err != nil {
		return 0, fmt.Errorf("read rdbuff: %w", err)
	}

	return v, nil
}

func (s *SWD) ReadCSW() (CSW, error) {
	v, err := s.ReadMemAP("CSW", regApCSW)
	if err != nil {
		return 0, fmt.Errorf("read MemAP: %w", err)
	}

	return CSW(v), nil
}

func (s *SWD) WriteCSW(csw CSW) error {
	return s.WriteMemAP("CSW", regApCSW, uint32(csw))
}

func (s *SWD) ReadIDR() (uint32, error) {
	return s.ReadMemAP("IDR", regIDR)
}

func (s *SWD) ReadBase() (uint32, error) {
	return s.ReadMemAP("BASE", regBase)
}

func (s *SWD) UpdateCSW(value, mask CSW) error {
	csw, err := s.ReadCSW()
	if err != nil {
		return fmt.Errorf("read CSW: %w", err)
	}

	csw &= ^mask
	csw |= value & mask

	return s.WriteCSW(csw)
}

func (s *SWD) WriteTAR(addr uint32) error {
	return s.WriteMemAP("TAR", regApTAR, addr)
}

func (s *SWD) WriteDRW(data uint32) error {
	return s.WriteMemAP("DRW", regApDRW, data)
}

func (s *SWD) ReadDRW() (uint32, error) {
	return s.ReadMemAP("DRW", regApDRW)
}

func (s *SWD) WriteRegister(addr uint32, data uint32) error {
	if err := s.WriteTAR(addr); err != nil {
		return fmt.Errorf("write TAR: %w", err)
	}

	if err := s.WriteDRW(data); err != nil {
		return fmt.Errorf("write DRW: %w", err)
	}

	if _, err := s.ReadRdBuff(); err != nil {
		return fmt.Errorf("read rdbuff: %w", err)
	}

	if _, err := s.ReadCtrlStat(); err != nil {
		return fmt.Errorf("read ctrlstat: %w", err)
	}

	return nil
}

func (s *SWD) ReadRegister(addr uint32) (uint32, error) {
	if err := s.WriteTAR(addr); err != nil {
		return 0, err
	}

	drw, err := s.ReadDRW()
	if err != nil {
		return 0, err
	}

	if _, err := s.ReadCtrlStat(); err != nil {
		return 0, fmt.Errorf("read ctrlstat: %w", err)
	}

	return drw, nil
}

func (s *SWD) UpdateRegisterBits(addr, mask, data uint32) error {
	v, err := s.ReadRegister(addr)
	if err != nil {
		return err
	}

	v &= ^mask
	v |= data & mask

	return s.WriteRegister(addr, v)
}

func (s *SWD) IDCode() (uint32, error) {
	return s.readTx("IDCODE", io.DebugPort, regIdCode)
}

func (s *SWD) ReadCtrlStat() (CtrlStat, error) {
	v, err := s.readTx("CTRL/STAT", io.DebugPort, regCtrlStat)
	if err != nil {
		return 0, err
	}

	return CtrlStat(v), nil
}

func (s *SWD) WriteCtrlStat(v CtrlStat) error {
	return s.writeTx("CTRL/STAT", io.DebugPort, regCtrlStat, uint32(v))
}

func (s *SWD) PowerOnReset() error {
	ctrlStat := CtrlStatDebugPowerUpRequest |
		CtrlStatSystemPowerUpRequest

	if err := s.WriteCtrlStat(ctrlStat); err != nil {
		// dummy read
		_, _ = s.ReadCtrlStat()

		return err
	}

	wanted := CtrlStatDebugPowerUpAck |
		CtrlStatSystemPowerUpAck

	for i := 0; i < 10; i++ {
		ctrlStat, err := s.ReadCtrlStat()
		if err != nil {
			return err
		}

		if ctrlStat&wanted == wanted {
			return nil
		}

		time.Sleep(time.Millisecond)
	}

	return ErrTimeout
}

func (s *SWD) Initialize() (uint32, error) {
	for i := 0; i < 100; i++ {
		if err := s.accessor.LineReset(); err != nil {
			return 0, fmt.Errorf("line reset: %w", err)
		}

		if err := s.accessor.LineReset(); err != nil {
			return 0, fmt.Errorf("line reset: %w", err)
		}

		id, err := s.IDCode()
		if err != nil {
			return 0, fmt.Errorf("idcode read: %w", err)
		}

		err = s.PowerOnReset()

		_ = s.Abort(AbortAllFlags())

		if err == nil {
			return id, nil
		}

		time.Sleep(time.Millisecond)
	}

	return 0, ErrTimeout
}

func (s *SWD) SetDebugger(d debug.Debugger) {
	s.debugger = d
}

func New(accessor io.Accessor) *SWD {
	return &SWD{
		accessor: accessor,
		debugger: &debug.NopDebugger{},
	}
}
