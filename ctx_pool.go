package dyntpl

import "sync"

// Context pool.
type CtxPool struct {
	p sync.Pool
}

// Default instance of context pool.
// You may use it directly as dyntpl.CP.Get()/Put() or using functions AcquireCtx()/ReleaseCtx().
var CP CtxPool

// Get context object from the pool or make new object if pool is empty.
func (p *CtxPool) Get() *Ctx {
	v := p.p.Get()
	if v != nil {
		if c, ok := v.(*Ctx); ok {
			return c
		}
	}
	return NewCtx()
}

// Put the object to the pool.
func (p *CtxPool) Put(ctx *Ctx) {
	ctx.Reset()
	p.p.Put(ctx)
}

// Get object from the default context pool.
func AcquireCtx() *Ctx {
	return CP.Get()
}

// Put object back to default pool.
func ReleaseCtx(ctx *Ctx) {
	CP.Put(ctx)
}
