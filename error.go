package dyntpl

import "errors"

var (
	ErrUnexpectedEOF = errors.New("unexpected end of file: control structure couldn't be closed")
	ErrBadCtl        = errors.New("bad control structure caught")
	ErrUnknownCtl    = errors.New("unknown ctl")

	ErrComplexCond   = errors.New("condition is too complex, use CondHelper instead")
	ErrSenselessCond = errors.New("comparison of two static args")

	ErrTplNotFound = errors.New("template not found")
	ErrInterrupt   = errors.New("tpl processing interrupted")
	ErrEmptyArg    = errors.New("empty input param")
	ErrModNoArgs   = errors.New("empty arguments list")

	ErrLoopParse     = errors.New("couldn't parse loop control structure")
	ErrWrongLoopLim  = errors.New("wrong count loop limit argument")
	ErrWrongLoopCond = errors.New("wrong loop condition operation")
	ErrWrongLoopOp   = errors.New("wrong loop operation")
	ErrBreakLoop     = errors.New("break loop")
	ErrContLoop      = errors.New("continue loop")
)
