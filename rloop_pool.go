package cbytetpl

import (
	"io"
	"sync"
)

type RangeLoopPool struct {
	p sync.Pool
}

var RLP RangeLoopPool

func (p *RangeLoopPool) Get(node Node, tpl *Tpl, ctx *Ctx, w io.Writer) *RangeLoop {
	v := p.p.Get()
	if v != nil {
		if rl, ok := v.(*RangeLoop); ok {
			rl.node = node
			rl.tpl = tpl
			rl.ctx = ctx
			rl.w = w
			return rl
		}
	}
	return NewRangeLoop(node, tpl, ctx, w)
}

func (p *RangeLoopPool) Put(rl *RangeLoop) {
	rl.ctx.Reset()
	rl.cntr = 0
	p.p.Put(rl)
}
