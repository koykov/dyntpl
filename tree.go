package dyntpl

import (
	"bytes"

	"github.com/koykov/bytebuf"
	"github.com/koykov/byteconv"
)

// Tree structure that represents parsed template as list of nodes with childrens.
type Tree struct {
	nodes []Node
	hsum  uint64
}

// Representation argument of modifier or helper.
type arg struct {
	name   []byte
	val    []byte
	static bool
}

var (
	hrQ  = []byte(`"`)
	hrQR = []byte(`&quot;`)

	ctlRepl = map[string]string{
		"\n":   "\\n",
		"\r":   "\\r",
		"\r\n": "\\r\\n",
		"\t":   "\\t",
	}
)

// HumanReadable builds human-readable view of the tree (currently in XML format).
func (t *Tree) HumanReadable() []byte {
	if len(t.nodes) == 0 {
		return nil
	}
	var buf bytebuf.Chain
	t.hrHelper(&buf, t.nodes, []byte("\t"), 0)
	return buf.Bytes()
}

// Internal human-readable helper.
func (t *Tree) hrHelper(buf *bytebuf.Chain, nodes []Node, indent []byte, depth int) {
	if depth == 0 {
		buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	} else {
		buf.Write(bytes.Repeat(indent, depth))
	}
	buf.WriteString("<nodes>\n")
	depth++
	for _, node := range nodes {
		buf.Write(bytes.Repeat(indent, depth))
		buf.WriteString(`<node type="`).WriteString(node.typ.String()).WriteByte('"')

		if len(node.prefix) > 0 {
			t.attrB(buf, "prefix", node.prefix)
		}
		if len(node.suffix) > 0 {
			t.attrB(buf, "suffix", node.suffix)
		}

		if len(node.ctxVar) > 0 && len(node.ctxSrc) > 0 {
			t.attrB(buf, "var", node.ctxVar)
			if len(node.ctxOK) > 0 {
				t.attrB(buf, "varOK", node.ctxOK)
			}
			t.attrB(buf, "src", node.ctxSrc)
			if len(node.ctxIns) > 0 {
				t.attrB(buf, "ins", node.ctxIns)
			}
		}

		if len(node.cntrVar) > 0 {
			t.attrB(buf, "var", node.cntrVar)
			if node.cntrInitF {
				t.attrI(buf, "val", node.cntrInit)
			} else {
				t.attrS(buf, "op", node.cntrOp.String())
				t.attrI(buf, "delta", node.cntrOpArg)
			}
		}

		if node.typ == TypeCond {
			if len(node.condL) > 0 {
				t.attrB(buf, "left", node.condL)
			}
			if node.condOp != 0 {
				t.attrS(buf, "op", node.condOp.String())
			}
			if len(node.condR) > 0 {
				t.attrB(buf, "right", node.condR)
			}
			if len(node.condHlp) > 0 {
				t.attrB(buf, "helper", node.condHlp)
				if node.condLC > lcNone {
					t.attrS(buf, "lc", node.condLC.String())
				}
				if len(node.condHlpArg) > 0 {
					for j, a := range node.condHlpArg {
						pfx := "arg"
						if a.static {
							pfx = "sarg"
						}
						buf.WriteByte(' ').
							WriteString(pfx).
							WriteInt(int64(j)).
							WriteString(`="`).
							Write(a.val).
							WriteByte('"')
					}
				}
			}
		}

		if node.typ == TypeCondOK {
			t.attrB(buf, "var", node.condOKL)
			t.attrB(buf, "varOK", node.condOKR)

			if len(node.condHlp) > 0 {
				t.attrB(buf, "helper", node.condHlp)
				if len(node.condHlpArg) > 0 {
					for j, a := range node.condHlpArg {
						pfx := "arg"
						if a.static {
							pfx = "sarg"
						}
						buf.WriteByte(' ').
							WriteString(pfx).
							WriteInt(int64(j)).
							WriteString(`="`).
							Write(a.val).
							WriteByte('"')
					}
				}
			}

			if len(node.condL) > 0 {
				t.attrB(buf, "left", node.condL)
			}
			if node.condOp != 0 {
				t.attrS(buf, "op", node.condOp.String())
			}
			if len(node.condR) > 0 {
				t.attrB(buf, "right", node.condR)
			}
		}

		if len(node.loopKey) > 0 {
			t.attrB(buf, "key", node.loopKey)
		}
		if len(node.loopVal) > 0 {
			t.attrB(buf, "val", node.loopVal)
		}
		if len(node.loopSrc) > 0 {
			t.attrB(buf, "src", node.loopSrc)
		}
		if len(node.loopCnt) > 0 {
			t.attrB(buf, "counter", node.loopCnt)
		}
		if node.loopCondOp != 0 {
			t.attrS(buf, "cond", node.loopCondOp.String())
		}
		if len(node.loopLim) > 0 {
			t.attrB(buf, "limit", node.loopLim)
		}
		if node.loopCntOp != 0 {
			t.attrS(buf, "op", node.loopCntOp.String())
		}
		if len(node.loopSep) > 0 {
			t.attrB(buf, "sep", node.loopSep)
		}
		if node.loopBrkD > 0 {
			t.attrI(buf, "brkD", node.loopBrkD)
		}

		if len(node.switchArg) > 0 {
			t.attrB(buf, "arg", node.switchArg)
		}
		if len(node.caseL) > 0 && node.caseOp != 0 && len(node.caseR) > 0 {
			t.attrB(buf, "left", node.caseL)
			t.attrS(buf, "op", node.caseOp.String())
			t.attrB(buf, "right", node.caseR)
		} else if len(node.caseL) > 0 {
			t.attrB(buf, "val", node.caseL)
		}
		if len(node.caseHlp) > 0 {
			t.attrB(buf, "helper", node.caseHlp)
			if len(node.caseHlpArg) > 0 {
				for j, a := range node.caseHlpArg {
					pfx := "arg"
					if a.static {
						pfx = "sarg"
					}
					buf.WriteByte(' ').
						WriteString(pfx).
						WriteInt(int64(j)).
						WriteString(`="`).
						Write(a.val).
						WriteByte('"')
				}
			}
		}

		if len(node.tpl) > 0 {
			for j, tpl := range node.tpl {
				buf.WriteByte(' ').
					WriteString("tpl").
					WriteInt(int64(j)).
					WriteString(`="`).
					Write(tpl).
					WriteByte('"')
			}
		}

		if len(node.loc) > 0 {
			t.attrB(buf, "val", node.loc)
		}

		if node.typ != TypeExit && node.typ != TypeBreak && node.typ != TypeLBreak && node.typ != TypeContinue && len(node.raw) > 0 {
			raw := string(node.raw)
			if repl, ok := ctlRepl[raw]; ok {
				raw = repl
			}
			t.attrS(buf, "val", raw)
		}

		if len(node.mod) > 0 || len(node.child) > 0 {
			buf.WriteByte('>')
		}

		if len(node.mod) > 0 {
			depth++
			buf.WriteByte('\n').Write(bytes.Repeat(indent, depth)).WriteString("<mods>\n")
			depth++
			for _, mod := range node.mod {
				buf.Write(bytes.Repeat(indent, depth)).WriteString(`<mod name="`).Write(mod.id).WriteByte('"')
				if len(mod.arg) > 0 {
					for j, a := range mod.arg {
						if len(a.name) > 0 {
							buf.WriteByte(' ').
								WriteString("key").
								WriteInt(int64(j)).
								WriteString(`="`).
								Write(a.name).
								WriteByte('"')
						}

						pfx := "arg"
						if a.static {
							pfx = "sarg"
						}
						v := a.val
						if bytes.Contains(v, hrQ) {
							v = bytes.ReplaceAll(v, hrQ, hrQR)
						}
						buf.WriteByte(' ').
							WriteString(pfx).
							WriteInt(int64(j)).
							WriteString(`="`).
							Write(v).
							WriteByte('"')
					}
				}
				buf.WriteString("/>\n")
			}
			depth--
			buf.Write(bytes.Repeat(indent, depth)).WriteString("</mods>\n")
			depth--
		}

		if len(node.child) > 0 {
			buf.WriteByte('\n')
			t.hrHelper(buf, node.child, indent, depth+1)
		}
		if len(node.mod) > 0 || len(node.child) > 0 {
			buf.Write(bytes.Repeat(indent, depth)).WriteString("</node>\n")
		} else {
			buf.WriteString("/>\n")
		}
	}
	depth--
	if depth > 0 {
		buf.Write(bytes.Repeat(indent, depth))
	}
	buf.WriteString("</nodes>\n")
}

func (t *Tree) attrB(buf *bytebuf.Chain, key string, p []byte) {
	buf.WriteByte(' ').
		WriteString(key).
		WriteString(`="`).
		Write(bytes.ReplaceAll(p, hrQ, hrQR)).
		WriteByte('"')
}

func (t *Tree) attrS(buf *bytebuf.Chain, key, s string) {
	t.attrB(buf, key, byteconv.S2B(s))
}

func (t *Tree) attrI(buf *bytebuf.Chain, key string, i int) {
	buf.WriteByte(' ').
		WriteString(key).
		WriteString(`="`).
		WriteInt(int64(i)).
		WriteByte('"')
}
