package systemcontrolblock

import "github.com/holoplot/go-swd/pkg/swd"

type SystemControlBlock struct {
	swd *swd.SWD
}

const (
	baseAddress = 0xe000ed00

	regAIRCR = baseAddress + 0xc
)

type AIRCR uint32

const (
	AIRCRVectReset       AIRCR = 1 << 0
	AIRCRVectClearActive AIRCR = 1 << 1
	AIRCRSysResetReq     AIRCR = 1 << 2
	AIRCRPriGroupShift         = 8
	AIRCRPriGroupMask    AIRCR = 0x7 << AIRCRPriGroupShift
	AIRCREndianessBig    AIRCR = 1 << 15
	AIRCREndianessLittle AIRCR = 0 << 15
	AIRCRVectKeyShift          = 16
	AIRCRVectKeyStat     AIRCR = 0x05fa << AIRCRVectKeyShift
)

func (cd *SystemControlBlock) ReadAIRCR() (AIRCR, error) {
	reg, err := cd.swd.ReadRegister(regAIRCR)
	if err != nil {
		return 0, err
	}

	return AIRCR(reg), nil
}

func (cd *SystemControlBlock) WriteAIRCR(aircr AIRCR) error {
	return cd.swd.WriteRegister(regAIRCR, uint32(aircr))
}

func (scb *SystemControlBlock) ResetSystem() error {
	return scb.WriteAIRCR(AIRCRVectKeyStat | AIRCRSysResetReq)
}

func New(swd *swd.SWD) *SystemControlBlock {
	return &SystemControlBlock{
		swd: swd,
	}
}
