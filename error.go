package dyntpl

import "errors"

var (
	ErrUnexpectedEOF = errors.New("unexpected end of file: control structure couldn't be closed")
	ErrUnknownCtl    = errors.New("unknown ctl")

	ErrSenselessCond   = errors.New("comparison of two static args")
	ErrCondHlpNotFound = errors.New("condition helper not found")

	ErrTplNotFound = errors.New("template not found")
	ErrInterrupt   = errors.New("tpl processing interrupted")
	ErrModNoArgs   = errors.New("empty arguments list")
	ErrModPoorArgs = errors.New("arguments list is too small")
	ErrModNoStr    = errors.New("argument is not string or bytes")
	ErrModEmptyStr = errors.New("argument is empty string")

	ErrWrongLoopLim  = errors.New("wrong count loop limit argument")
	ErrWrongLoopCond = errors.New("wrong loop condition operation")
	ErrWrongLoopOp   = errors.New("wrong loop operation")
	ErrBreakLoop     = errors.New("break loop")
	ErrLBreakLoop    = errors.New("lazybreak loop")
	ErrContLoop      = errors.New("continue loop")
)
