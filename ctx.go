package cbytetpl

import "github.com/koykov/cbytealg"

type Ctx struct {
	vars  []ctxVar
	ssbuf []string
}

type ctxVar struct {
	key string
	val interface{}
	ins Inspector
}

func NewCtx() *Ctx {
	ctx := Ctx{
		vars:  make([]ctxVar, 0),
		ssbuf: make([]string, 0),
	}
	return &ctx
}

func (c *Ctx) Set(key string, val interface{}) {
	c.SetWithInspector(key, val, ReflectInspector{})
}

func (c *Ctx) SetWithInspector(key string, val interface{}, ins Inspector) {
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
			return v.ins.Get(v.val, c.ssbuf[1:]...)
		}
	}

	return nil
}
