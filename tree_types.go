package dyntpl

// Type of the node.
type Type int

// Op is a type of the operation in conditions and loops.
type Op int

type lc int

const (
	TypeRaw       Type = 0
	TypeTpl       Type = 1
	TypeCond      Type = 2
	TypeCondOK    Type = 3
	TypeCondTrue  Type = 4
	TypeCondFalse Type = 5
	TypeLoopRange Type = 6
	TypeLoopCount Type = 7
	TypeBreak     Type = 8
	TypeLBreak    Type = 9
	TypeContinue  Type = 10
	TypeCtx       Type = 11
	TypeCounter   Type = 12
	TypeSwitch    Type = 13
	TypeCase      Type = 14
	TypeDefault   Type = 15
	TypeDiv       Type = 16
	TypeJsonQ     Type = 17
	TypeEndJsonQ  Type = 18
	TypeHtmlE     Type = 19
	TypeEndHtmlE  Type = 20
	TypeUrlEnc    Type = 21
	TypeEndUrlEnc Type = 22
	TypeInclude   Type = 23
	TypeExit      Type = 99

	// Must be in sync with inspector.Op type.
	OpUnk Op = 0
	OpEq  Op = 1
	OpNq  Op = 2
	OpGt  Op = 3
	OpGtq Op = 4
	OpLt  Op = 5
	OpLtq Op = 6
	OpInc Op = 7
	OpDec Op = 8

	lcNone lc = 0
	lcLen  lc = 1
	lcCap  lc = 2
)

// String view of the type.
func (typ Type) String() string {
	switch typ {
	case TypeRaw:
		return "raw"
	case TypeTpl:
		return "tpl"
	case TypeCond:
		return "cond"
	case TypeCondOK:
		return "condOK"
	case TypeCondTrue:
		return "true"
	case TypeCondFalse:
		return "false"
	case TypeLoopRange:
		return "rloop"
	case TypeLoopCount:
		return "cloop"
	case TypeBreak:
		return "break"
	case TypeLBreak:
		return "lazybreak"
	case TypeContinue:
		return "cont"
	case TypeCtx:
		return "ctx"
	case TypeCounter:
		return "cntr"
	case TypeSwitch:
		return "switch"
	case TypeCase:
		return "case"
	case TypeDefault:
		return "def"
	case TypeDiv:
		return "div"
	case TypeInclude:
		return "inc"
	case TypeExit:
		return "exit"
	default:
		return "unk"
	}
}

// String view of the operation.
func (o Op) String() string {
	switch o {
	case OpEq:
		return "=="
	case OpNq:
		return "!="
	case OpGt:
		return ">"
	case OpGtq:
		return ">="
	case OpLt:
		return "<"
	case OpLtq:
		return "<="
	case OpInc:
		return "++"
	case OpDec:
		return "--"
	default:
		return "unk"
	}
}

// Swap inverts itself.
func (o Op) Swap() Op {
	switch o {
	case OpGt:
		return OpLt
	case OpGtq:
		return OpLtq
	case OpLt:
		return OpGt
	case OpLtq:
		return OpGtq
	default:
		return o
	}
}

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
