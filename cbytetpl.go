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
		for _, bcFn := range byteConvFnRegistry {
			ctx.bbuf = ctx.bbuf[:0]
			ctx.bbuf, err = bcFn(ctx.bbuf, raw)
			if err == nil && len(ctx.bbuf) > 0 {
				_, err = t.w.Write(ctx.bbuf)
				break
			}
		}
	default:
		err = ErrUnknownCtl
	}
	if err == ErrUnknownType {
		return
	}
	return
}

func init() {
	RegisterByteConvFn(byteConvBytes)
	RegisterByteConvFn(byteConvStr)
	RegisterByteConvFn(byteConvBool)
	RegisterByteConvFn(byteConvInt)
	RegisterByteConvFn(byteConvUint)
	RegisterByteConvFn(byteConvFloat)
}
