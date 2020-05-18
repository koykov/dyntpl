package dyntpl

import (
	"bytes"
)

type Tree struct {
	nodes []Node
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
		if node.typ != TypeExit && node.typ != TypeBreak && node.typ != TypeContinue {
			buf.WriteByte(':')
			buf.WriteByte(' ')
			buf.Write(node.raw)
		}
		if len(node.prefix) > 0 {
			buf.WriteString(" pfx ")
			buf.Write(node.prefix)
		}
		if len(node.suffix) > 0 {
			buf.WriteString(" sfx ")
			buf.Write(node.suffix)
		}

		if len(node.ctxVar) > 0 && len(node.ctxSrc) > 0 {
			buf.WriteString("var ")
			buf.Write(node.ctxVar)
			buf.WriteString(" src ")
			buf.Write(node.ctxSrc)
			if len(node.ctxIns) > 0 {
				buf.WriteString(" ins ")
				buf.Write(node.ctxIns)
			}
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

		if len(node.loopKey) > 0 {
			buf.WriteString("key ")
			buf.Write(node.loopKey)
			buf.WriteByte(' ')
		}
		if len(node.loopVal) > 0 {
			buf.WriteString("val ")
			buf.Write(node.loopVal)
		}
		if len(node.loopSrc) > 0 {
			buf.WriteString(" src ")
			buf.Write(node.loopSrc)
		}
		if len(node.loopCnt) > 0 {
			buf.WriteString("cnt ")
			buf.Write(node.loopCnt)
		}
		if node.loopCondOp != 0 {
			buf.WriteString(" cond ")
			buf.WriteString(node.loopCondOp.String())
		}
		if len(node.loopLim) > 0 {
			buf.WriteString(" lim ")
			buf.Write(node.loopLim)
		}
		if node.loopCntOp != 0 {
			buf.WriteString(" op ")
			buf.WriteString(node.loopCntOp.String())
		}
		if len(node.loopSep) > 0 {
			buf.WriteString(" sep ")
			buf.Write(node.loopSep)
		}

		if len(node.switchArg) > 0 {
			buf.WriteString("arg ")
			buf.Write(node.switchArg)
		}
		if len(node.caseL) > 0 && node.caseOp != 0 && len(node.caseR) > 0 {
			buf.WriteString("left ")
			buf.Write(node.caseL)
			buf.WriteString(" op ")
			buf.WriteString(node.caseOp.String())
			buf.WriteString(" right ")
			buf.Write(node.caseR)
		} else if len(node.caseL) > 0 {
			buf.WriteString("val ")
			buf.Write(node.caseL)
		}

		if len(node.mod) > 0 {
			buf.WriteString(" mod")
			for i, mod := range node.mod {
				if i > 0 {
					buf.WriteByte(',')
				}
				buf.WriteByte(' ')
				buf.Write(mod.id)
				if len(mod.arg) > 0 {
					buf.WriteByte('(')
					for j, a := range mod.arg {
						if j > 0 {
							buf.WriteByte(',')
							buf.WriteByte(' ')
						}
						if a.static {
							buf.WriteByte('"')
							buf.Write(a.val)
							buf.WriteByte('"')
						} else {
							buf.Write(a.val)
						}
					}
					buf.WriteByte(')')
				}
			}
		}

		buf.WriteByte('\n')
		if len(node.child) > 0 {
			t.hrHelper(buf, node.child, indent, depth+1)
		}
	}
}
