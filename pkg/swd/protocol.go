package swd

import "github.com/holoplot/go-swd/pkg/io"

const (
	regIdCode     io.Address = 0x0
	regAbort      io.Address = 0x0
	regCtrlStat   io.Address = 0x4
	regResend     io.Address = 0x8
	regSelect     io.Address = 0x8
	regReadBuffer io.Address = 0xc

	regApCSW io.Address = 0x0
	regApTAR io.Address = 0x4
	regApDRW io.Address = 0xc

	regBase io.Address = 0xf8
	regIDR  io.Address = 0xfc
)

type AbortFlags uint32

const (
	AbortDAP                AbortFlags = 1 << 0
	AbortStickyCmpClear     AbortFlags = 1 << 1
	AbortStickyErrClear     AbortFlags = 1 << 2
	AbortWdErrorClear       AbortFlags = 1 << 3
	AbortStickyOverrunClear AbortFlags = 1 << 4
)

func AbortAllFlags() AbortFlags {
	return AbortStickyCmpClear |
		AbortStickyErrClear |
		AbortWdErrorClear |
		AbortStickyOverrunClear
}

type CSW uint32

const (
	CSWAutoIncrementOff    CSW = 0x0 << 4
	CSWAutoIncrementSingle CSW = 0x1 << 4
	CSWAutoIncrementPacked CSW = 0x2 << 4
	CSWAutoIncrementMask   CSW = 0x3 << 4

	CSWDeviceEnable       CSW = 1 << 6
	CSWTransferInProgress CSW = 1 << 7

	CSWModeShift     = 8
	CSWModeMask  CSW = 0xf << CSWModeShift

	CSWBusAccessProtectionShift = 24
	CSWBusAccessProtectionMask  = 0x7f << CSWBusAccessProtectionShift

	CSWDebugSoftwareEnable CSW = 1 << 31

	CSWSize8bit  CSW = 0x0
	CSWSize16bit CSW = 0x1
	CSWSize32bit CSW = 0x2
	CSWSizeMask  CSW = 0x7
)

type CtrlStat uint32

const (
	CtrlStatOverrunDetect        CtrlStat = 1 << 0
	CtrlStatStickyOverrunDetect  CtrlStat = 1 << 1
	CtrlStatStickyCmp            CtrlStat = 1 << 4
	CtrlStatStickyErr            CtrlStat = 1 << 5
	CtrlStatReadOk               CtrlStat = 1 << 6
	CtrlStatWriteDataError       CtrlStat = 1 << 7
	CtrlStatDebugResetRequest    CtrlStat = 1 << 26
	CtrlStatDebugResetAck        CtrlStat = 1 << 27
	CtrlStatDebugPowerUpRequest  CtrlStat = 1 << 28
	CtrlStatDebugPowerUpAck      CtrlStat = 1 << 29
	CtrlStatSystemPowerUpRequest CtrlStat = 1 << 30
	CtrlStatSystemPowerUpAck     CtrlStat = 1 << 31

	CtrlStatTransferModeShift          = 2
	CtrlStatTransferModeMask  CtrlStat = 0x3 << CtrlStatTransferModeShift

	CtrlStatMaskLaneShift          = 8
	CtrlStatMaskLaneMask  CtrlStat = 0xf << CtrlStatMaskLaneShift

	CtrlStatTransactionCounterShift          = 12
	CtrlStatTransactionCounterMask  CtrlStat = 0xfff << CtrlStatTransactionCounterShift
)

func (ctrlStat *CtrlStat) TransferMode() uint8 {
	return uint8((*ctrlStat & CtrlStatTransferModeMask) >> CtrlStatTransferModeShift)
}

func (ctrlStat *CtrlStat) SetTransferMode(mode uint8) {
	*ctrlStat &= ^CtrlStatTransferModeMask
	*ctrlStat |= CtrlStat(mode) << CtrlStatTransferModeShift
}

func (ctrlStat *CtrlStat) MaskLane() uint8 {
	return uint8((*ctrlStat & CtrlStatMaskLaneMask) >> CtrlStatMaskLaneShift)
}

func (ctrlStat *CtrlStat) SetMaskLane(lane uint8) {
	*ctrlStat &= ^CtrlStatMaskLaneMask
	*ctrlStat |= CtrlStat(lane) << CtrlStatMaskLaneShift
}

func (ctrlStat *CtrlStat) TransactionCounter() uint16 {
	return uint16((*ctrlStat & CtrlStatTransactionCounterMask) >> CtrlStatTransactionCounterShift)
}

func (ctrlStat *CtrlStat) SetTransactionCounter(counter uint16) {
	*ctrlStat &= ^CtrlStatTransactionCounterMask
	*ctrlStat |= CtrlStat(counter) << CtrlStatTransactionCounterShift
}
