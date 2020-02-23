package cbytetpl

import (
	"bytes"
)

type Type int

const (
	TypeRaw Type = iota
	TypeTpl
	TypeCondition
	TypeLoop
	TypeCtx
	TypeSwitch
	TypeCase
	TypeDefault
)

func (typ Type) String() string {
	switch typ {
	case TypeRaw:
		return "raw"
	case TypeTpl:
		return "tpl"
	case TypeCondition:
		return "cnd"
	case TypeLoop:
		return "lop"
	case TypeCtx:
		return "ctx"
	case TypeSwitch:
		return "swt"
	case TypeCase:
		return "cas"
	case TypeDefault:
		return "def"
	default:
		return "unk"
	}
}

type Tree struct {
	nodes []Node
}

type Node struct {
	typ   Type
	raw   []byte
	child []Node
}

func (t *Tree) humanReadable() []byte {
	if len(t.nodes) == 0 {
		return nil
	}
	var buf bytes.Buffer
	t.hrHelper(&buf, t.nodes, []byte("  "), 0)
	return buf.Bytes()
}

func (t *Tree) hrHelper(buf *bytes.Buffer, nodes []Node, indent []byte, depth int) {
	for _, node := range nodes {
		buf.Write(bytes.Repeat(indent, depth))
		buf.WriteString(node.typ.String())
		buf.WriteByte(':')
		buf.WriteByte(' ')
		buf.Write(node.raw)
		buf.WriteByte('\n')
		if len(node.child) > 0 {
			t.hrHelper(buf, node.child, indent, depth+1)
		}
	}
}
