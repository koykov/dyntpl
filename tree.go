package dyntpl

import (
	"bytes"

	"github.com/koykov/bytebuf"
	"github.com/koykov/byteconv"
)

// Tree structure that represents parsed template as list of nodes with childrens.
type Tree struct {
	nodes []node
	hsum  uint64
}

// Representation argument of modifier or helper.
type arg struct {
	name   []byte
	val    []byte
	static bool
	global bool
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
func (t *Tree) hrHelper(buf *bytebuf.Chain, nodes []node, indent []byte, depth int) {
	if depth == 0 {
		buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n")
	} else {
		buf.Write(bytes.Repeat(indent, depth))
	}
	buf.WriteString("<nodes>\n")
	depth++
	for _, n := range nodes {
		buf.Write(bytes.Repeat(indent, depth))
		buf.WriteString(`<node type="`).WriteString(n.typ.String()).WriteByte('"')

		if len(n.prefix) > 0 {
			t.attrB(buf, "prefix", n.prefix)
		}
		if len(n.suffix) > 0 {
			t.attrB(buf, "suffix", n.suffix)
		}

		if len(n.ctxVar) > 0 && len(n.ctxSrc) > 0 {
			t.attrB(buf, "var", n.ctxVar)
			if len(n.ctxOK) > 0 {
				t.attrB(buf, "varOK", n.ctxOK)
			}
			t.attrB(buf, "src", n.ctxSrc)
			if len(n.ctxIns) > 0 {
				t.attrB(buf, "ins", n.ctxIns)
			}
		}

		if len(n.cntrVar) > 0 {
			t.attrB(buf, "var", n.cntrVar)
			if n.cntrInitF {
				t.attrI(buf, "val", n.cntrInit)
			} else {
				t.attrS(buf, "op", n.cntrOp.String())
				t.attrI(buf, "delta", n.cntrOpArg)
			}
		}

		if n.typ == typeCond {
			if len(n.condL) > 0 {
				t.attrB(buf, "left", n.condL)
			}
			if n.condOp != 0 {
				t.attrS(buf, "op", n.condOp.String())
			}
			if len(n.condR) > 0 {
				t.attrB(buf, "right", n.condR)
			}
			if len(n.condHlp) > 0 {
				t.attrB(buf, "helper", n.condHlp)
				if n.condLC > lcNone {
					t.attrS(buf, "lc", n.condLC.String())
				}
				if len(n.condHlpArg) > 0 {
					for j, a := range n.condHlpArg {
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

		if n.typ == typeCondOK {
			t.attrB(buf, "var", n.condOKL)
			t.attrB(buf, "varOK", n.condOKR)

			if len(n.condHlp) > 0 {
				t.attrB(buf, "helper", n.condHlp)
				if len(n.condHlpArg) > 0 {
					for j, a := range n.condHlpArg {
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

			if len(n.condL) > 0 {
				t.attrB(buf, "left", n.condL)
			}
			if n.condOp != 0 {
				t.attrS(buf, "op", n.condOp.String())
			}
			if len(n.condR) > 0 {
				t.attrB(buf, "right", n.condR)
			}
		}

		if len(n.loopKey) > 0 {
			t.attrB(buf, "key", n.loopKey)
		}
		if len(n.loopVal) > 0 {
			t.attrB(buf, "val", n.loopVal)
		}
		if len(n.loopSrc) > 0 {
			t.attrB(buf, "src", n.loopSrc)
		}
		if len(n.loopCnt) > 0 {
			t.attrB(buf, "counter", n.loopCnt)
		}
		if n.loopCondOp != 0 {
			t.attrS(buf, "cond", n.loopCondOp.String())
		}
		if len(n.loopLim) > 0 {
			t.attrB(buf, "limit", n.loopLim)
		}
		if n.loopCntOp != 0 {
			t.attrS(buf, "op", n.loopCntOp.String())
		}
		if len(n.loopSep) > 0 {
			t.attrB(buf, "sep", n.loopSep)
		}
		if n.loopBrkD > 0 {
			t.attrI(buf, "brkD", n.loopBrkD)
		}

		if len(n.switchArg) > 0 {
			t.attrB(buf, "arg", n.switchArg)
		}
		if len(n.caseL) > 0 && n.caseOp != 0 && len(n.caseR) > 0 {
			t.attrB(buf, "left", n.caseL)
			t.attrS(buf, "op", n.caseOp.String())
			t.attrB(buf, "right", n.caseR)
		} else if len(n.caseL) > 0 {
			t.attrB(buf, "val", n.caseL)
		}
		if len(n.caseHlp) > 0 {
			t.attrB(buf, "helper", n.caseHlp)
			if len(n.caseHlpArg) > 0 {
				for j, a := range n.caseHlpArg {
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

		if len(n.tpl) > 0 {
			for j, tpl := range n.tpl {
				buf.WriteByte(' ').
					WriteString("tpl").
					WriteInt(int64(j)).
					WriteString(`="`).
					Write(tpl).
					WriteByte('"')
			}
		}

		if len(n.loc) > 0 {
			t.attrB(buf, "val", n.loc)
		}

		if n.typ != typeExit && n.typ != typeBreak && n.typ != typeLBreak && n.typ != typeContinue && len(n.raw) > 0 {
			raw := string(n.raw)
			if repl, ok := ctlRepl[raw]; ok {
				raw = repl
			}
			t.attrS(buf, "val", raw)
		}

		if n.noesc {
			t.attrS(buf, "noesc", "true")
		}

		if len(n.mod) > 0 || len(n.child) > 0 {
			buf.WriteByte('>')
		}

		if len(n.mod) > 0 {
			depth++
			buf.WriteByte('\n').Write(bytes.Repeat(indent, depth)).WriteString("<mods>\n")
			depth++
			for _, mod := range n.mod {
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

		if len(n.child) > 0 {
			buf.WriteByte('\n')
			t.hrHelper(buf, n.child, indent, depth+1)
		}
		if len(n.mod) > 0 || len(n.child) > 0 {
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
