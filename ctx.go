package cbytetpl

import (
	"bytes"
	"io"
	"strconv"

	"github.com/koykov/cbytealg"
	"github.com/koykov/fastconv"
	"github.com/koykov/inspector"
)

type Ctx struct {
	vars  []ctxVar
	ln    int
	ssbuf []string
	bbuf  []byte
	bbuf1 []byte
	cbuf  bool
	ibuf  int64
	chQB  bool
	buf   interface{}
	rl    *RangeLoop
	Err   error
}

type ctxVar struct {
	key string
	val interface{}
	buf []byte
	ins inspector.Inspector
}

var (
	sqL = []byte("[")
	sqR = []byte("]")
	dot = []byte(".")
)

func NewCtx() *Ctx {
	ctx := Ctx{
		vars:  make([]ctxVar, 0),
		ssbuf: make([]string, 0),
		bbuf:  make([]byte, 0),
	}
	return &ctx
}

func (c *Ctx) Set(key string, val interface{}, ins inspector.Inspector) {
	for i := 0; i < c.ln; i++ {
		if c.vars[i].key == key {
			c.vars[i].val = val
			c.vars[i].ins = ins
			return
		}
	}
	if c.ln < len(c.vars) {
		c.vars[c.ln].val = val
		c.vars[c.ln].ins = ins
	} else {
		c.vars = append(c.vars, ctxVar{
			key: key,
			val: val,
			ins: ins,
		})
	}
	c.ln++
}

func (c *Ctx) SetStatic(key string, val interface{}) {
	ins, err := inspector.GetInspector("static")
	if err != nil {
		c.Err = err
		return
	}
	c.Set(key, val, ins)
}

func (c *Ctx) SetBytes(key string, val []byte) {
	ins, err := inspector.GetInspector("static")
	if err != nil {
		c.Err = err
		return
	}
	for i := 0; i < c.ln; i++ {
		if c.vars[i].key == key {
			c.vars[i].buf = append(c.vars[i].buf[:0], val...)
			c.vars[i].ins = ins
			return
		}
	}
	if c.ln < len(c.vars) {
		c.vars[c.ln].key = key
		c.vars[c.ln].buf = append(c.vars[c.ln].buf[:0], val...)
		c.vars[c.ln].ins = ins
	} else {
		v := ctxVar{
			key: key,
			buf: val,
			ins: ins,
		}
		c.vars = append(c.vars, v)
	}
	c.ln++
}

func (c *Ctx) Get(path string) interface{} {
	return c.get(fastconv.S2B(path))
}

func (c *Ctx) Reset() {
	c.Err = nil
	c.buf = nil
	c.chQB = false
	c.ln = 0
	c.ssbuf = c.ssbuf[:0]
	c.bbuf = c.bbuf[:0]
	c.bbuf1 = c.bbuf1[:0]
}

func (c *Ctx) get(path []byte) interface{} {
	if c.chQB {
		path = c.replaceQB(path)
	}

	c.ssbuf = c.ssbuf[:0]
	c.ssbuf = cbytealg.AppendSplitStr(c.ssbuf, fastconv.B2S(path), ".", -1)
	if len(c.ssbuf) == 0 {
		return nil
	}

	for i, v := range c.vars {
		if i == c.ln {
			break
		}
		if v.key == c.ssbuf[0] {
			if v.val == nil && v.buf != nil {
				c.bbuf = append(c.bbuf[:0], v.buf...)
				c.buf = &c.bbuf
				return c.buf
			}
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

func (c *Ctx) rloop(path []byte, node Node, tpl *Tpl, w io.Writer) {
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

func (c *Ctx) cloop(node Node, tpl *Tpl, w io.Writer) {
	var (
		cnt, lim  int64
		cntr      int
		allowIter bool
	)
	cnt = c.cloopRange(node.loopCntStatic, node.loopCntInit)
	if c.Err != nil {
		return
	}
	lim = c.cloopRange(node.loopLimStatic, node.loopLim)
	if c.Err != nil {
		return
	}
	c.ibuf = cnt
	allowIter = false
	cntr = 0
	for {
		switch node.loopCondOp {
		case OpLt:
			allowIter = c.ibuf < lim
		case OpLtq:
			allowIter = c.ibuf <= lim
		case OpGt:
			allowIter = c.ibuf > lim
		case OpGtq:
			allowIter = c.ibuf >= lim
		case OpEq:
			allowIter = c.ibuf == lim
		case OpNq:
			allowIter = c.ibuf != lim
		default:
			c.Err = ErrWrongLoopCond
			break
		}
		if !allowIter {
			break
		}

		c.Set(fastconv.B2S(node.loopCnt), &c.ibuf, &inspector.StaticInspector{})

		if cntr > 0 && len(node.loopSep) > 0 {
			_, _ = w.Write(node.loopSep)
		}
		cntr++
		c.chQB = true
		for _, ch := range node.child {
			_ = tpl.renderNode(w, ch, c)
		}
		c.chQB = false

		switch node.loopCntOp {
		case OpInc:
			c.ibuf++
		case OpDec:
			c.ibuf--
		default:
			c.Err = ErrWrongLoopOp
			break
		}
	}
	return
}

func (c *Ctx) cloopRange(static bool, b []byte) (r int64) {
	if static {
		r, c.Err = strconv.ParseInt(fastconv.B2S(b), 0, 0)
		if c.Err != nil {
			return
		}
	} else {
		var ok bool
		raw := c.get(b)
		if c.Err != nil {
			return
		}
		r, ok = lim2int(raw)
		if !ok {
			c.Err = ErrWrongLoopLim
			return
		}
	}
	return
}

func (c *Ctx) replaceQB(path []byte) []byte {
	sqLi := bytes.Index(path, sqL)
	sqRi := bytes.Index(path, sqR)
	if sqLi != -1 && sqRi != -1 && sqLi < sqRi && sqRi < len(path) {
		c.bbuf = append(c.bbuf[:0], path[0:sqLi]...)
		c.bbuf = append(c.bbuf, dot...)
		c.chQB = false
		c.buf = c.get(path[sqLi+1 : sqRi])
		if c.buf != nil {
			c.bbuf1, c.Err = cbytealg.AnyToBytes(c.bbuf1[:0], c.buf)
			if c.Err != nil {
				c.chQB = true
				return nil
			}
			c.bbuf = append(c.bbuf, c.bbuf1...)
		}
		c.chQB = true
		c.bbuf = append(c.bbuf, path[sqRi+1:]...)
		path = c.bbuf
	}
	return path
}
