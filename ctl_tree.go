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
		return "cond"
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
	child  []Node
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
		if len(node.prefix) > 0 {
			buf.WriteString(" pfx ")
			buf.Write(node.prefix)
		}
		if len(node.suffix) > 0 {
			buf.WriteString(" sfx ")
			buf.Write(node.suffix)
		}
		buf.WriteByte('\n')
		if len(node.child) > 0 {
			t.hrHelper(buf, node.child, indent, depth+1)
		}
	}
}

func (t *Tree) addNode(node Node) {
	t.nodes = append(t.nodes, node)
}

func (t *Tree) addRaw(raw []byte) {
	if len(raw) == 0 {
		return
	}
	t.nodes = append(t.nodes, Node{typ: TypeRaw, raw: raw})
}

func (t *Tree) addTpl(tpl []byte) {
	t.nodes = append(t.nodes, Node{typ: TypeTpl, raw: tpl})
}
