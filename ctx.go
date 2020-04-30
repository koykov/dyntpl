package cbytetpl

import (
	"github.com/koykov/cbytealg"
	"github.com/koykov/inspector"
)

type Ctx struct {
	vars  []ctxVar
	ssbuf []string
	buf   interface{}
	Err   error
}

type ctxVar struct {
	key string
	val interface{}
	ins inspector.Inspector
}

func NewCtx() *Ctx {
	ctx := Ctx{
		vars:  make([]ctxVar, 0),
		ssbuf: make([]string, 0),
	}
	return &ctx
}

func (c *Ctx) Set(key string, val interface{}, ins inspector.Inspector) {
	c.vars = append(c.vars, ctxVar{
		key: key,
		val: val,
		ins: ins,
	})
}

func (c *Ctx) Get(path string) interface{} {
	c.ssbuf = c.ssbuf[:0]
	c.ssbuf = cbytealg.AppendSplitStr(c.ssbuf, path, ".", -1)
	if len(c.ssbuf) == 0 {
		return nil
	}

	for _, v := range c.vars {
		if v.key == c.ssbuf[0] {
			c.Err = v.ins.GetTo(v.val, &c.buf, c.ssbuf[1:]...)
			if c.Err != nil {
				return nil
			}
			return c.buf
		}
	}

	return nil
}

func (c *Ctx) Reset() {
	c.Err = nil
	c.buf = nil
	c.vars = c.vars[:0]
	c.ssbuf = c.ssbuf[:0]
}
