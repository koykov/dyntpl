package cbytetpl

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/koykov/cbytealg"
)

var (
	empty    []byte
	noFmt    = []byte(" \t\n")
	ctlOpen  = []byte("{%")
	ctlClose = []byte("%}")

	reCutComments = regexp.MustCompile(`\t*{#[^#]*#}\n*`)
	reCutFmt      = regexp.MustCompile(`\n+\t*\s*`)
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

	cl, cs int
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
			tree.nodes = append(tree.nodes, Node{typ: TypeRaw, raw: p.tpl[o:]})
			o = i
			break
		}
		if inCtl {
			e := cbytealg.IndexAt(p.tpl, ctlClose, i)
			if e < 0 {
				return nil, ErrUnexpectedEOF
			}
			e += 2
			tree.nodes = append(tree.nodes, Node{typ: TypeTpl, raw: p.tpl[o:e]})
			o, i = e, e
			inCtl = false
		} else {
			tree.nodes = append(tree.nodes, Node{typ: TypeRaw, raw: p.tpl[o:i]})
			o = i
			inCtl = true
		}
	}
	return tree, nil
}
