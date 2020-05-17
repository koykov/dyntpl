package cbytetpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/koykov/cbytealg"
	"github.com/koykov/fastconv"
)

const (
	targetCond = iota
	targetLoop
	targetSwitch
)

type Parser struct {
	keepFmt bool
	tpl     []byte

	// Counters of conditions, loops and switches.
	cc, cl, cs int
}

type target map[int]int

var (
	empty      []byte
	space      = []byte(" ")
	comma      = []byte(",")
	uscore     = []byte("_")
	vline      = []byte("|")
	quotes     = []byte("\"'`")
	noFmt      = []byte(" \t\n")
	ctlExit    = []byte("exit")
	ctlOpen    = []byte("{%")
	ctlClose   = []byte("%}")
	ctlTrim    = []byte("{}% ")
	ctlTrimAll = []byte("{}%= ")
	ctxStatic  = []byte("static")
	condElse   = []byte("else")
	condEnd    = []byte("endif")
	loopEnd    = []byte("endfor")
	loopBrk    = []byte("break")
	loopCnt    = []byte("continue")
	swDefault  = []byte("default")
	swEnd      = []byte("endswitch")

	opEq  = []byte("==")
	opNq  = []byte("!=")
	opGt  = []byte(">")
	opGtq = []byte(">=")
	opLt  = []byte("<")
	opLtq = []byte("<=")
	opInc = []byte("++")
	opDec = []byte("--")

	reCutComments = regexp.MustCompile(`\t*{#[^#]*#}\n*`)
	reCutFmt      = regexp.MustCompile(`\n+\t*\s*`)

	reTplPS = regexp.MustCompile(`^=\s*(.*) (?:prefix|pfx) (.*) (?:suffix|sfx) (.*)`)
	reTplP  = regexp.MustCompile(`^=\s*(.*) (?:prefix|pfx) (.*)`)
	reTplS  = regexp.MustCompile(`^=\s*(.*) (?:suffix|sfx) (.*)`)
	reTpl   = regexp.MustCompile(`^= (.*)`)
	reMod   = regexp.MustCompile(`([^(]+)\(*([^)]*)\)*`)

	reCtx = regexp.MustCompile(`ctx (\w+)\s*=\s*([\w.]+)\s*[as]*\s*(\w*)`)

	reCond        = regexp.MustCompile(`if .*`)
	reCondExpr    = regexp.MustCompile(`if (.*)(==|!=|>=|<=|>|<)(.*)`)
	reCondComplex = regexp.MustCompile(`if .*&&|\|\||\(|\).*`)

	reLoop      = regexp.MustCompile(`for .*`)
	reLoopRange = regexp.MustCompile(`for ([^:]+)\s*:*=\s*range\s*([^\s]*)\s*(?:separator|sep)*\s*(.*)`)
	reLoopCount = regexp.MustCompile(`for (\w*)\s*:*=\s*(\w+)\s*;\s*\w+\s*(<|<=|>|>=|!=)+\s*([^;]+)\s*;\s*\w*(--|\+\+)+\s*(?:separator|sep)*\s*(.*)`)

	reSwitch     = regexp.MustCompile(`^switch\s*(.*)`)
	reSwitchCase = regexp.MustCompile(`case ([^<=>!]+)([<=>!]{2})*(.*)`)
)

func Parse(tpl []byte, keepFmt bool) (tree *Tree, err error) {
	p := &Parser{
		tpl:     tpl,
		keepFmt: keepFmt,
	}
	p.cutComments()
	p.cutFmt()

	tree = &Tree{}
	target := newTarget(p)
	tree.nodes, _, err = p.parseTpl(tree.nodes, 0, target)
	return
}

func ParseFile(fileName string, keepFmt bool) (tree *Tree, err error) {
	_, err = os.Stat(fileName)
	if os.IsNotExist(err) {
		return
	}
	var raw []byte
	raw, err = ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file %s", fileName)
	}
	return Parse(raw, keepFmt)
}

func (p *Parser) cutComments() {
	p.tpl = reCutComments.ReplaceAll(p.tpl, empty)
}

func (p *Parser) cutFmt() {
	if p.keepFmt {
		return
	}
	p.tpl = reCutFmt.ReplaceAll(p.tpl, empty)
	p.tpl = cbytealg.Trim(p.tpl, noFmt)
}

func (p *Parser) parseTpl(nodes []Node, offset int, target *target) ([]Node, int, error) {
	var (
		up  bool
		err error
	)

	o, i := offset, offset
	inCtl := false
	for !target.reached(p) || target.eqZero() {
		i = cbytealg.IndexAt(p.tpl, ctlOpen, i)
		if i < 0 {
			if inCtl {
				return nodes, o, ErrUnexpectedEOF
			}
			nodes = addRaw(nodes, p.tpl[o:])
			o = len(p.tpl)
			break
		}
		if inCtl {
			e := cbytealg.IndexAt(p.tpl, ctlClose, i)
			if e < 0 {
				return nodes, o, ErrUnexpectedEOF
			}
			e += 2
			node := Node{}
			nodes, e, up, err = p.processCtl(nodes, &node, p.tpl[o:e], o)
			if err != nil {
				return nodes, o, err
			}
			o, i = e, e
			inCtl = false
			if up {
				break
			}
		} else {
			nodes = addRaw(nodes, p.tpl[o:i])
			o = i
			inCtl = true
		}
	}
	return nodes, o, nil
}

func (p *Parser) processCtl(nodes []Node, root *Node, ctl []byte, pos int) ([]Node, int, bool, error) {
	var (
		offset int
		up     bool
		err    error
	)

	up = false
	t := cbytealg.Trim(ctl, ctlTrim)
	// Check tpl structure
	if reTplPS.Match(t) || reTplP.Match(t) || reTplS.Match(t) || reTpl.Match(t) {
		root.typ = TypeTpl
		if m := reTplPS.FindSubmatch(t); m != nil {
			root.raw, root.mod = p.extractMods(m[1])
			root.prefix = m[2]
			root.suffix = m[3]
		} else if m := reTplP.FindSubmatch(t); m != nil {
			root.raw, root.mod = p.extractMods(m[1])
			root.prefix = m[2]
		} else if m := reTplS.FindSubmatch(t); m != nil {
			root.raw, root.mod = p.extractMods(m[1])
			root.suffix = m[2]
		} else {
			root.raw, root.mod = p.extractMods(cbytealg.Trim(t, ctlTrimAll))
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check context structure.
	if reCtx.Match(t) {
		root.typ = TypeCtx
		if m := reCtx.FindSubmatch(t); m != nil {
			root.ctxVar = m[1]
			root.ctxSrc = m[2]
			root.ctxSrcStatic = isStatic(m[2])
			if len(m) > 3 && len(m[3]) > 0 {
				root.ctxIns = m[3]
			} else {
				root.ctxIns = ctxStatic
			}
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check condition structure.
	if reCond.Match(t) {
		if reCondComplex.Match(t) {
			return nodes, pos, up, ErrComplexCond
		}
		target := newTarget(p)
		p.cc++

		subNodes := make([]Node, 0)
		subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), target)
		split := splitNodes(subNodes)

		root.typ = TypeCond
		root.condL, root.condR, root.condStaticL, root.condStaticR, root.condOp = p.parseCondExpr(t)
		if len(split) > 0 {
			nodeTrue := Node{typ: TypeCondTrue, child: split[0]}
			root.child = append(root.child, nodeTrue)
		}
		if len(split) > 1 {
			nodeFalse := Node{typ: TypeCondFalse, child: split[1]}
			root.child = append(root.child, nodeFalse)
		}

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check condition divider.
	if bytes.Equal(t, condElse) {
		root.typ = TypeDiv
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check condition end.
	if bytes.Equal(t, condEnd) {
		p.cc--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}

	// Check loop structure.
	if reLoop.Match(t) {
		if m := reLoopRange.FindSubmatch(t); m != nil {
			root.typ = TypeLoopRange
			if bytes.Contains(m[1], comma) {
				kv := bytes.Split(m[1], comma)
				root.loopKey = cbytealg.Trim(kv[0], space)
				if bytes.Equal(root.loopKey, uscore) {
					root.loopKey = nil
				}
				root.loopVal = cbytealg.Trim(kv[1], space)
			} else {
				root.loopKey = cbytealg.Trim(m[1], space)
			}
			root.loopSrc = m[2]
			if len(m) > 2 {
				root.loopSep = m[3]
			}
		} else if m := reLoopCount.FindSubmatch(t); m != nil {
			root.typ = TypeLoopCount
			root.loopCnt = m[1]
			root.loopCntInit = m[2]
			root.loopCntStatic = isStatic(m[2])
			root.loopCondOp = p.parseOp(m[3])
			root.loopLim = m[4]
			root.loopLimStatic = isStatic(m[4])
			root.loopCntOp = p.parseOp(m[5])
			if len(m) > 5 {
				root.loopSep = m[6]
			}
		} else {
			return nodes, 0, up, ErrLoopParse
		}

		target := newTarget(p)
		p.cl++

		root.child = make([]Node, 0)
		root.child, offset, err = p.parseTpl(root.child, pos+len(ctl), target)

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check loop end.
	if bytes.Equal(t, loopEnd) {
		p.cl--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}
	// Check loop break.
	if bytes.Equal(t, loopBrk) {
		root.typ = TypeBreak
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check loop continue.
	if bytes.Equal(t, loopCnt) {
		root.typ = TypeContinue
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check switch structure.
	if m := reSwitch.FindSubmatch(t); m != nil {
		target := newTarget(p)
		p.cs++

		root.typ = TypeSwitch
		if len(m) > 0 {
			root.switchArg = m[1]
		}
		root.child = make([]Node, 0)
		root.child, offset, err = p.parseTpl(root.child, pos+len(ctl), target)
		root.child = rollupSwitchNodes(root.child)

		nodes = addNode(nodes, *root)
		return nodes, offset, up, err
	}
	// Check switch's case.
	if reSwitchCase.Match(t) {
		root.typ = TypeCase
		root.caseL, root.caseR, root.caseStaticL, root.caseStaticR, root.caseOp = p.parseCaseExpr(t)
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch's default.
	if bytes.Equal(t, swDefault) {
		root.typ = TypeDefault
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}
	// Check switch end.
	if bytes.Equal(t, swEnd) {
		p.cs--
		offset = pos + len(ctl)
		up = true
		return nodes, offset, up, err
	}

	// Check tpl interrupt.
	if bytes.Equal(t, ctlExit) {
		root.typ = TypeExit
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	return nodes, 0, up, ErrBadCtl
}

func (p *Parser) parseCondExpr(expr []byte) (l, r []byte, sl, sr bool, op Op) {
	if m := reCondExpr.FindSubmatch(expr); m != nil {
		l = cbytealg.Trim(m[1], space)
		r = cbytealg.Trim(m[3], space)
		sl = isStatic(l)
		sr = isStatic(r)
		op = p.parseOp(m[2])
	}
	return
}

func (p *Parser) parseCaseExpr(expr []byte) (l, r []byte, sl, sr bool, op Op) {
	if m := reSwitchCase.FindSubmatch(expr); m != nil {
		l = cbytealg.Trim(m[1], space)
		sl = isStatic(l)
		if len(m) > 1 {
			op = p.parseOp(m[2])
			r = cbytealg.Trim(m[3], space)
			sr = isStatic(r)
		}
	}
	return
}

func (p *Parser) parseOp(src []byte) Op {
	var op Op
	switch {
	case bytes.Equal(src, opEq):
		op = OpEq
	case bytes.Equal(src, opNq):
		op = OpNq
	case bytes.Equal(src, opGt):
		op = OpGt
	case bytes.Equal(src, opGtq):
		op = OpGtq
	case bytes.Equal(src, opLt):
		op = OpLt
	case bytes.Equal(src, opLtq):
		op = OpLtq
	case bytes.Equal(src, opInc):
		op = OpInc
	case bytes.Equal(src, opDec):
		op = OpDec
	default:
		op = OpUnk
	}
	return op
}

func (p *Parser) extractMods(t []byte) ([]byte, []mod) {
	if bytes.Contains(t, vline) {
		chunks := bytes.Split(t, vline)
		mods := make([]mod, 0, len(chunks)-1)
		for i := 1; i < len(chunks); i++ {
			if m := reMod.FindSubmatch(chunks[i]); m != nil {
				fn := GetModFn(fastconv.B2S(m[1]))
				if fn == nil {
					continue
				}
				arg := make([]*modArg, 0)
				if len(m) > 1 {
					args := bytes.Split(m[2], comma)
					for _, a := range args {
						a = cbytealg.Trim(a, space)
						arg = append(arg, &modArg{
							val:    cbytealg.Trim(a, quotes),
							static: isStatic(a),
						})
					}
				}
				mods = append(mods, mod{
					id:  m[1],
					fn:  fn,
					arg: arg,
				})
			}
		}
		return chunks[0], mods
	} else {
		return t, nil
	}
}

func newTarget(p *Parser) *target {
	return &target{
		targetCond:   p.cc,
		targetLoop:   p.cl,
		targetSwitch: p.cs,
	}
}

func (t *target) reached(p *Parser) bool {
	return (*t)[targetCond] == p.cc &&
		(*t)[targetLoop] == p.cl &&
		(*t)[targetSwitch] == p.cs
}

func (t *target) eqZero() bool {
	return (*t)[targetCond] == 0 &&
		(*t)[targetLoop] == 0 &&
		(*t)[targetSwitch] == 0
}
