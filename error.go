package cbytetpl

import "errors"

var (
	ErrUnexpectedEOF = errors.New("unexpected end of file: control structure couldn't be closed")
	ErrBadCtl        = errors.New("bad control structure caught")
	ErrUnknownCtl    = errors.New("unknown ctl")
	ErrCondComplex   = errors.New("condition is too complex, use CondHelper instead")
	ErrLoopParse     = errors.New("couldn't parse loop control structure")
	ErrTplNotFound   = errors.New("template not found")
	ErrEmptyArg      = errors.New("empty input param")
	ErrSenselessCond = errors.New("comparison of two static args")
	ErrWrongLoopLim  = errors.New("wrong count loop limit argument")
	ErrWrongLoopCond = errors.New("wrong loop condition operation")
	ErrWrongLoopOp   = errors.New("wrong loop operation")

	ErrModNoArgs = errors.New("empty arguments list")
)
