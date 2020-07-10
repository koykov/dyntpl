package dyntpl

// Type of the node.
type Type int

// Type of the operation in conditions and loops.
type Op int

const (
	// Known types of nodes.
	TypeRaw       Type = 0
	TypeTpl       Type = 1
	TypeCond      Type = 2
	TypeCondTrue  Type = 3
	TypeCondFalse Type = 4
	TypeLoopRange Type = 5
	TypeLoopCount Type = 6
	TypeBreak     Type = 7
	TypeContinue  Type = 8
	TypeCtx       Type = 9
	TypeCounter   Type = 10
	TypeSwitch    Type = 11
	TypeCase      Type = 12
	TypeDefault   Type = 13
	TypeDiv       Type = 14
	TypeJsonQ     Type = 15
	TypeEndJsonQ  Type = 16
	TypeHtmlE     Type = 17
	TypeEndHtmlE  Type = 18
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
	case TypeExit:
		return "exit"
	default:
		return "unk"
	}
}

// String view of the opertion.
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

// Invert operation.
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
