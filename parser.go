package cbytetpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/koykov/cbytealg"
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
	noSpace    = []byte(" ")
	noFmt      = []byte(" \t\n")
	ctlOpen    = []byte("{%")
	ctlClose   = []byte("%}")
	ctlTrim    = []byte("{}% ")
	ctlTrimAll = []byte("{}%= ")
	condElse   = []byte(`else`)
	condEnd    = []byte(`endif`)

	opEq  = []byte("==")
	opNq  = []byte("!=")
	opGt  = []byte(">")
	opGtq = []byte(">=")
	opLt  = []byte("<")
	opLtq = []byte("<=")

	reCutComments = regexp.MustCompile(`\t*{#[^#]*#}\n*`)
	reCutFmt      = regexp.MustCompile(`\n+\t*\s*`)

	reTplPS = regexp.MustCompile(`^=\s*(.*) (?:prefix|pfx) (.*) (?:suffix|sfx) (.*)`)
	reTplP  = regexp.MustCompile(`^=\s*(.*) (?:prefix|pfx) (.*)`)
	reTplS  = regexp.MustCompile(`^=\s*(.*) (?:suffix|sfx) (.*)`)
	reTpl   = regexp.MustCompile(`^= (.*)`)

	reCond        = regexp.MustCompile(`if .*`)
	reCondExpr    = regexp.MustCompile(`if (.*)(==|!=|>|>=|<|<=)(.*)`)
	reCondComplex = regexp.MustCompile(`if .*&&|\|\||\(|\).*`)
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
			root.raw = m[1]
			root.prefix = m[2]
			root.suffix = m[3]
		} else if m := reTplP.FindSubmatch(t); m != nil {
			root.raw = m[1]
			root.prefix = m[2]
		} else if m := reTplS.FindSubmatch(t); m != nil {
			root.raw = m[1]
			root.suffix = m[2]
		} else {
			root.raw = cbytealg.Trim(t, ctlTrimAll)
		}
		nodes = addNode(nodes, *root)
		offset = pos + len(ctl)
		return nodes, offset, up, err
	}

	// Check condition structure.
	if reCond.Match(t) {
		if reCondComplex.Match(t) {
			return nodes, pos, up, ErrCondComplex
		}
		target := newTarget(p)
		p.cc++

		subNodes := make([]Node, 0)
		subNodes, offset, err = p.parseTpl(subNodes, pos+len(ctl), target)
		split := splitNodes(subNodes)

		root.typ = TypeCond
		root.condL, root.condR, root.condOp = p.parseCondExpr(t)
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

	return nodes, 0, up, ErrBadCtl
}

func (p *Parser) parseCondExpr(expr []byte) (l, r []byte, op CondOp) {
	if m := reCondExpr.FindSubmatch(expr); m != nil {
		l = cbytealg.Trim(m[1], noSpace)
		r = cbytealg.Trim(m[3], noSpace)
		switch {
		case bytes.Equal(m[2], opEq):
			op = CondOpEq
		case bytes.Equal(m[2], opNq):
			op = CondOpNq
		case bytes.Equal(m[2], opGt):
			op = CondOpGt
		case bytes.Equal(m[2], opGtq):
			op = CondOpGtq
		case bytes.Equal(m[2], opLt):
			op = CondOpLt
		case bytes.Equal(m[2], opLtq):
			op = CondOpLtq
		default:
			op = CondOpUnk
		}
	}
	return
}

func newTarget(p *Parser) *target {
	return &target{targetCond: p.cc, targetLoop: p.cl, targetSwitch: p.cs}
}

func (t *target) reached(p *Parser) bool {
	return (*t)[targetCond] == p.cc && (*t)[targetLoop] == p.cl && (*t)[targetSwitch] == p.cs
}

func (t *target) eqZero() bool {
	return (*t)[targetCond] == 0 && (*t)[targetLoop] == 0 && (*t)[targetSwitch] == 0
}
