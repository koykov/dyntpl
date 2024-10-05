package dyntpl

import "sync"

// pool is a context pool.
type pool struct {
	p sync.Pool
}

// cp is a default instance of context pool.
var cp pool

// AcquireCtx gets object from the default context pool.
func AcquireCtx() *Ctx {
	v := cp.p.Get()
	if v != nil {
		if c, ok := v.(*Ctx); ok {
			return c
		}
	}
	return NewCtx()
}

// ReleaseCtx puts object back to default pool.
func ReleaseCtx(ctx *Ctx) {
	ctx.Reset()
	cp.p.Put(ctx)
}
