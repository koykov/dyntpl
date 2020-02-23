package cbytetpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/koykov/cbytealg"
)

var (
	empty      []byte
	noFmt      = []byte(" \t\n")
	ctlOpen    = []byte("{%")
	ctlClose   = []byte("%}")
	ctlTrim    = []byte("{}% ")
	ctlTrimAll = []byte("{}%= ")

	reCutComments = regexp.MustCompile(`\t*{#[^#]*#}\n*`)
	reCutFmt      = regexp.MustCompile(`\n+\t*\s*`)

	reTplPS = regexp.MustCompile(`=\s*(.*) prefix (.*) suffix (.*)`)
	reTplP  = regexp.MustCompile(`=\s*(.*) prefix (.*)`)
	reTplS  = regexp.MustCompile(`=\s*(.*) suffix (.*)`)
	reTpl   = regexp.MustCompile(`= (.*)`)
)

func Parse(tpl []byte, keepFmt bool) (tree *Tree, err error) {
	p := &Parser{
		tpl:     tpl,
		keepFmt: keepFmt,
	}
	p.cutComments()
	p.cutFmt()
	return p.parseTpl()
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

type Parser struct {
	keepFmt bool
	tpl     []byte

	// Counters of conditions, loops and switches.
	cc, cl, cs int
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

func (p *Parser) parseTpl() (*Tree, error) {
	tree := &Tree{}
	o, i := 0, 0
	inCtl := false
	for {
		i = cbytealg.IndexAt(p.tpl, ctlOpen, i)
		if i < 0 {
			if inCtl {
				return nil, ErrUnexpectedEOF
			}
			tree.addRaw(p.tpl[o:])
			o = i
			break
		}
		if inCtl {
			e := cbytealg.IndexAt(p.tpl, ctlClose, i)
			if e < 0 {
				return nil, ErrUnexpectedEOF
			}
			e += 2
			node := Node{}
			e, err := p.processCtl(tree, &node, p.tpl[o:e], o)
			if err != nil {
				return nil, err
			}
			o, i = e, e
			inCtl = false
		} else {
			tree.addRaw(p.tpl[o:i])
			o = i
			inCtl = true
		}
	}
	return tree, nil
}

func (p *Parser) processCtl(tree *Tree, root *Node, ctl []byte, pos int) (offset int, err error) {
	t := cbytealg.Trim(ctl, ctlTrim)
	// Check tpl control
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
		tree.addNode(*root)
		offset = pos + len(ctl)
		return
	}

	return 0, ErrBadCtl
}
