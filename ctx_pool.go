package cbytetpl

import "sync"

type CtxPool struct {
	p sync.Pool
}

var CP CtxPool

func (p *CtxPool) Get() *Ctx {
	v := p.p.Get()
	if v != nil {
		if c, ok := v.(*Ctx); ok {
			return c
		}
	}
	return NewCtx()
}

func (p *CtxPool) Put(ctx *Ctx) {
	ctx.Reset()
	p.p.Put(ctx)
}
