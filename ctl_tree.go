package cbytetpl

import (
	"bytes"
)

type Type int
type CondOp int

const (
	TypeRaw Type = iota
	TypeTpl
	TypeCond
	TypeCondTrue
	TypeCondFalse
	TypeLoop
	TypeCtx
	TypeSwitch
	TypeCase
	TypeDefault
	TypeDiv

	CondOpUnk CondOp = iota
	CondOpEq
	CondOpNq
	CondOpGt
	CondOpGtq
	CondOpLt
	CondOpLtq
)

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
	case TypeLoop:
		return "loop"
	case TypeCtx:
		return "ctx"
	case TypeSwitch:
		return "switch"
	case TypeCase:
		return "case"
	case TypeDefault:
		return "def"
	case TypeDiv:
		return "div"
	default:
		return "unk"
	}
}

func (c CondOp) String() string {
	switch c {
	case CondOpEq:
		return "=="
	case CondOpNq:
		return "!="
	case CondOpGt:
		return ">"
	case CondOpGtq:
		return ">="
	case CondOpLt:
		return "<"
	case CondOpLtq:
		return "<="
	default:
		return "unk"
	}
}

type Tree struct {
	nodes []Node
}

type Node struct {
	typ    Type
	raw    []byte
	prefix []byte
	suffix []byte

	condL  []byte
	condR  []byte
	condOp CondOp

	child []Node
}

func (t *Tree) humanReadable() []byte {
	if len(t.nodes) == 0 {
		return nil
	}
	var buf bytes.Buffer
	t.hrHelper(&buf, t.nodes, []byte("\t"), 0)
	return buf.Bytes()
}

func (t *Tree) hrHelper(buf *bytes.Buffer, nodes []Node, indent []byte, depth int) {
	for _, node := range nodes {
		buf.Write(bytes.Repeat(indent, depth))
		buf.WriteString(node.typ.String())
		buf.WriteByte(':')
		buf.WriteByte(' ')
		buf.Write(node.raw)
		if len(node.prefix) > 0 {
			buf.WriteString(" pfx ")
			buf.Write(node.prefix)
		}
		if len(node.suffix) > 0 {
			buf.WriteString(" sfx ")
			buf.Write(node.suffix)
		}
		if len(node.condL) > 0 {
			buf.WriteString("left ")
			buf.Write(node.condL)
		}
		if node.condOp != 0 {
			buf.WriteString(" op ")
			buf.WriteString(node.condOp.String())
		}
		if len(node.condR) > 0 {
			buf.WriteString(" right ")
			buf.Write(node.condR)
		}
		buf.WriteByte('\n')
		if len(node.child) > 0 {
			t.hrHelper(buf, node.child, indent, depth+1)
		}
	}
}

func addNode(nodes []Node, node Node) []Node {
	nodes = append(nodes, node)
	return nodes
}

func addRaw(nodes []Node, raw []byte) []Node {
	if len(raw) == 0 {
		return nodes
	}
	nodes = append(nodes, Node{typ: TypeRaw, raw: raw})
	return nodes
}

func addTpl(nodes []Node, tpl []byte) []Node {
	nodes = append(nodes, Node{typ: TypeTpl, raw: tpl})
	return nodes
}

func splitNodes(nodes []Node) [][]Node {
	if len(nodes) == 0 {
		return nil
	}
	split := make([][]Node, 0)
	var o int
	for i, node := range nodes {
		if node.typ == TypeDiv {
			split = append(split, nodes[o:i])
			o = i + 1
		}
	}
	if o < len(nodes) {
		split = append(split, nodes[o:])
	}
	return split
}
