package swd

type ControlStatusRegister struct {
	CDBGPWRUPACK bool
	CDBGPWRUPREQ bool
	CDBGRSTACK   bool
	CDBGRSTREQ   bool
	RAZ          uint8
}
