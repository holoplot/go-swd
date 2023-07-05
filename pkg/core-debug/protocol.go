package coredebug

// https://developer.arm.com/documentation/ddi0337/e/core-debug/core-debug-registers

// Debug Halting Control and Status Register
type DHCSR uint32

const (
	// Read-write
	DHCSRDebugKey  DHCSR = 0xa05f << 16
	DHCSRCDebugEn  DHCSR = 1 << 0
	DHCSRCHalt     DHCSR = 1 << 1
	DHCSRCStep     DHCSR = 1 << 2
	DHCSRCMaskInts DHCSR = 1 << 3
	DHCSRCSnapAll  DHCSR = 1 << 5

	// Read-only
	DHCSRSRegReady     DHCSR = 1 << 16
	DHCSRSHalt         DHCSR = 1 << 17
	DHCSRSSleep        DHCSR = 1 << 18
	DHCSRSLockup       DHCSR = 1 << 19
	DHCSRSRetireStatus DHCSR = 1 << 24
	DHCSRSResetStatus  DHCSR = 1 << 25
)

// Debug Core Register Selector Register
type DCRSR uint32

const (
	RegSelMask DCRSR = 0x1f
	RegWnR     DCRSR = 1 << 16
)

// Debug Core Register Data Register
type DCRDR uint32

// Debug Exception and Monitor Control Register
type DEMCR uint32

const (
	DEMCRVcCoreReset       DEMCR = 1 << 0
	DEMCRVcMmErr           DEMCR = 1 << 4
	DEMCRVcNoCoproessorErr DEMCR = 1 << 5
	DEMCRVcCheckErr        DEMCR = 1 << 6
	DEMCRVcStateErr        DEMCR = 1 << 7
	DEMCRVcBusErr          DEMCR = 1 << 8
	DEMCRVcIntErr          DEMCR = 1 << 9
	DEMCRVcHardFaultErr    DEMCR = 1 << 10
	DEMCRMonitoringEnabled DEMCR = 1 << 16
	DEMCRMonitoringPend    DEMCR = 1 << 17
	DEMCRMonitoringStep    DEMCR = 1 << 18
	DEMCRMonitoringReq     DEMCR = 1 << 19
	DEMCREnableTrace       DEMCR = 1 << 24
)

const (
	baseAddress = 0xe000edf0

	regDHCSR = baseAddress + 0x0
	regDCRSR = baseAddress + 0x4
	regDCRDR = baseAddress + 0x8
	regDEMCR = baseAddress + 0xc
)
