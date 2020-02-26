package cbytetpl

import "errors"

var (
	ErrUnexpectedEOF = errors.New("unexpected end of file: control structure couldn't be closed")
	ErrBadCtl        = errors.New("bad control structure caught")
	ErrCondComplex   = errors.New("condition is too complex, use CondHelper instead")
	ErrLoopParse     = errors.New("couldn't parse loop control structure")
)
