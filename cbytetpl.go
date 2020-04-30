package cbytetpl

import (
	"bytes"
	"io"
	"sync"

	"github.com/koykov/fastconv"
)

type Tpl struct {
	Id   string
	tree *Tree
	w    io.Writer
}

var (
	mux         sync.Mutex
	tplRegistry = map[string]*Tpl{}
)

func RegisterTpl(id string, tree *Tree) {
	tpl := Tpl{
		Id:   id,
		tree: tree,
	}
	mux.Lock()
	tplRegistry[id] = &tpl
	mux.Unlock()
}

func Render(id string, ctx *Ctx) ([]byte, error) {
	buf := bytes.Buffer{}
	err := RenderTo(&buf, id, ctx)
	return buf.Bytes(), err
}

func RenderTo(w io.Writer, id string, ctx *Ctx) (err error) {
	mux.Lock()
	tpl, ok := tplRegistry[id]
	mux.Unlock()
	if !ok {
		err = ErrTplNotFound
		return
	}
	tpl.w = w
	for _, node := range tpl.tree.nodes {
		err = tpl.renderNode(&node, ctx)
		if err != nil {
			return
		}
	}

	return
}

func (t *Tpl) renderNode(node *Node, ctx *Ctx) (err error) {
	switch node.typ {
	case TypeRaw:
		_, err = t.w.Write(node.raw)
	case TypeTpl:
		raw := ctx.Get(fastconv.B2S(node.raw))
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
		switch raw.(type) {
		// Bytes case.
		case []byte:
			_, err = t.w.Write(raw.([]byte))
		case *[]byte:
			_, err = t.w.Write(*raw.(*[]byte))
			if err != nil {
				return
			}
		// String case.
		case string:
			_, err = t.w.Write(fastconv.S2B(raw.(string)))
		case *string:
			_, err = t.w.Write(fastconv.S2B(*raw.(*string)))
			if err != nil {
				return
			}
		// All other cases.
		default:
			//
		}
	default:
		err = ErrUnknownCtl
	}
	return
}
