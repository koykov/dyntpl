package cbytetpl

import (
	"io"

	"github.com/koykov/cbytealg"
	"github.com/koykov/fastconv"
	"github.com/koykov/inspector"
)

type Ctx struct {
	vars  []ctxVar
	ssbuf []string
	bbuf  []byte
	bbuf1 []byte
	cbuf  bool
	buf   interface{}
	rl    *RangeLoop
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
		bbuf:  make([]byte, 0),
	}
	return &ctx
}

func (c *Ctx) Set(key string, val interface{}, ins inspector.Inspector) {
	for i := range c.vars {
		if c.vars[i].key == key {
			c.vars[i].val = val
			c.vars[i].ins = ins
			return
		}
	}
	c.vars = append(c.vars, ctxVar{
		key: key,
		val: val,
		ins: ins,
	})
}

func (c *Ctx) Get(path string) interface{} {
	return c.get(fastconv.S2B(path))
}

func (c *Ctx) get(path []byte) interface{} {
	c.ssbuf = c.ssbuf[:0]
	c.ssbuf = cbytealg.AppendSplitStr(c.ssbuf, fastconv.B2S(path), ".", -1)
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

func (c *Ctx) cmp(path []byte, cond Op, right []byte) bool {
	c.ssbuf = c.ssbuf[:0]
	c.ssbuf = cbytealg.AppendSplitStr(c.ssbuf, fastconv.B2S(path), ".", -1)
	if len(c.ssbuf) == 0 {
		return false
	}

	for _, v := range c.vars {
		if v.key == c.ssbuf[0] {
			c.Err = v.ins.Cmp(v.val, inspector.Op(cond), fastconv.B2S(right), &c.cbuf, c.ssbuf[1:]...)
			if c.Err != nil {
				return false
			}
			return c.cbuf
		}
	}

	return false
}

func (c *Ctx) loop(path []byte, node Node, tpl *Tpl, w io.Writer) {
	c.ssbuf = c.ssbuf[:0]
	c.ssbuf = cbytealg.AppendSplitStr(c.ssbuf, fastconv.B2S(path), ".", -1)
	if len(c.ssbuf) == 0 {
		return
	}
	for _, v := range c.vars {
		if v.key == c.ssbuf[0] {
			if c.rl == nil {
				c.rl = NewRangeLoop(node, tpl, c, w)
			} else {
				c.rl.cntr = 0
				c.rl.node = node
				c.rl.tpl = tpl
				c.rl.w = w
			}
			c.Err = v.ins.Loop(v.val, c.rl, &c.bbuf1, c.ssbuf[1:]...)
			return
		}
	}
}

func (c *Ctx) Reset() {
	c.Err = nil
	c.buf = nil
	c.vars = c.vars[:0]
	c.ssbuf = c.ssbuf[:0]
	c.bbuf = c.bbuf[:0]
	c.bbuf1 = c.bbuf1[:0]
}
