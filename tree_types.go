package dyntpl

// rtype represents a runtime type of the node.
type rtype int

const (
	typeRaw rtype = iota
	typeTpl
	typeCond
	typeCondOK
	typeCondTrue
	typeCondFalse
	typeLoopRange
	typeLoopCount
	typeBreak
	typeLBreak
	typeContinue
	typeCtx
	typeCounter
	typeSwitch
	typeCase
	typeDefault
	typeDiv
	typeJsonQ
	typeEndJsonQ
	typeHtmlE
	typeEndHtmlE
	typeUrlEnc
	typeEndUrlEnc
	typeInclude
	typeExit
)

// String view of the type.
func (typ rtype) String() string {
	switch typ {
	case typeRaw:
		return "raw"
	case typeTpl:
		return "tpl"
	case typeCond:
		return "cond"
	case typeCondOK:
		return "condOK"
	case typeCondTrue:
		return "true"
	case typeCondFalse:
		return "false"
	case typeLoopRange:
		return "rloop"
	case typeLoopCount:
		return "cloop"
	case typeBreak:
		return "break"
	case typeLBreak:
		return "lazybreak"
	case typeContinue:
		return "cont"
	case typeCtx:
		return "ctx"
	case typeCounter:
		return "cntr"
	case typeSwitch:
		return "switch"
	case typeCase:
		return "case"
	case typeDefault:
		return "def"
	case typeDiv:
		return "div"
	case typeInclude:
		return "inc"
	case typeExit:
		return "exit"
	default:
		return "unk"
	}
}

// op represents a type of the operation in conditions and loops.
type op int

// Must be in sync with inspector.Op type.
const (
	opUnk op = iota
	opEq
	opNq
	opGt
	opGtq
	opLt
	opLtq
	opInc
	opDec
)

// Swap inverts itself.
func (o op) Swap() op {
	switch o {
	case opGt:
		return opLt
	case opGtq:
		return opLtq
	case opLt:
		return opGt
	case opLtq:
		return opGtq
	default:
		return o
	}
}

// String view of the operation.
func (o op) String() string {
	switch o {
	case opEq:
		return "=="
	case opNq:
		return "!="
	case opGt:
		return ">"
	case opGtq:
		return ">="
	case opLt:
		return "<"
	case opLtq:
		return "<="
	case opInc:
		return "++"
	case opDec:
		return "--"
	default:
		return "unk"
	}
}

// lc represents len/cap type.
type lc int

const (
	lcNone lc = iota
	lcLen
	lcCap
)

func (lc lc) String() string {
	switch lc {
	case lcLen:
		return "len"
	case lcCap:
		return "cap"
	case lcNone:
		fallthrough
	default:
		return ""
	}
}
