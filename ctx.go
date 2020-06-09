package dyntpl

import (
	"bytes"
	"io"
	"strconv"

	"github.com/koykov/any2bytes"
	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
	"github.com/koykov/inspector"
)

// Context object. Contains list of variables available to inspect.
// In addition has buffers to help develop new modifiers without allocations.
type Ctx struct {
	// List of context variables and list len.
	vars []ctxVar
	ln   int
	// Check square brackets flag.
	chQB bool
	// Internal buffers.
	buf  []byte
	bufS []string
	bufX interface{}
	bufA []interface{}
	// Range loop helper.
	rl *RangeLoop

	// External buffers to use in modifier and condition helpers.
	Buf, Buf1, Buf2 ByteBuf

	BufB bool
	BufI int64
	BufU uint64
	BufF float64

	Err error
}

// Context variable object.
type ctxVar struct {
	key string
	val interface{}
	// Byte buffer need for special cases when value is a byte slice.
	buf []byte
	ins inspector.Inspector
}

var (
	// Byte constants.
	qbL = []byte("[")
	qbR = []byte("]")
	dot = []byte(".")
)

func NewCtx() *Ctx {
	ctx := Ctx{
		vars: make([]ctxVar, 0),
		bufS: make([]string, 0),
		Buf:  make(ByteBuf, 0),
		Buf1: make(ByteBuf, 0),
		Buf2: make(ByteBuf, 0),
		buf:  make([]byte, 0),
		bufA: make([]interface{}, 0),
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
		c.vars[c.ln].key = key
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
	c.bufX = nil
	c.chQB = false
	c.ln = 0
	c.bufS = c.bufS[:0]
	c.Buf.Reset()
	c.Buf1.Reset()
	c.Buf2.Reset()
	c.buf = c.buf[:0]
	c.bufA = c.bufA[:0]
	if c.rl != nil {
		c.rl.Reset()
	}
}

func (c *Ctx) get(path []byte) interface{} {
	if c.chQB {
		path = c.replaceQB(path)
	}

	c.bufS = c.bufS[:0]
	c.bufS = bytealg.AppendSplitStr(c.bufS, fastconv.B2S(path), ".", -1)
	if len(c.bufS) == 0 {
		return nil
	}

	for i, v := range c.vars {
		if i == c.ln {
			break
		}
		if v.key == c.bufS[0] {
			if v.val == nil && v.buf != nil {
				c.Buf.Write(v.buf)
				c.bufX = &c.Buf
				return c.bufX
			}
			c.Err = v.ins.GetTo(v.val, &c.bufX, c.bufS[1:]...)
			if c.Err != nil {
				return nil
			}
			return c.bufX
		}
	}

	return nil
}

func (c *Ctx) cmp(path []byte, cond Op, right []byte) bool {
	c.bufS = c.bufS[:0]
	c.bufS = bytealg.AppendSplitStr(c.bufS, fastconv.B2S(path), ".", -1)
	if len(c.bufS) == 0 {
		return false
	}

	for i, v := range c.vars {
		if i == c.ln {
			break
		}
		if v.key == c.bufS[0] {
			c.Err = v.ins.Cmp(v.val, inspector.Op(cond), fastconv.B2S(right), &c.BufB, c.bufS[1:]...)
			if c.Err != nil {
				return false
			}
			return c.BufB
		}
	}

	return false
}

func (c *Ctx) rloop(path []byte, node Node, tpl *Tpl, w io.Writer) {
	c.bufS = c.bufS[:0]
	c.bufS = bytealg.AppendSplitStr(c.bufS, fastconv.B2S(path), ".", -1)
	if len(c.bufS) == 0 {
		return
	}
	for i, v := range c.vars {
		if i == c.ln {
			break
		}
		if v.key == c.bufS[0] {
			var rl *RangeLoop
			if c.rl == nil {
				c.rl = NewRangeLoop(node, tpl, c, w)
				rl = c.rl
			} else {
				crl := c.rl
				for {
					if crl.stat == rlFree {
						rl = crl
						break
					}
					if crl.stat != rlFree {
						if crl.next != nil {
							crl = crl.next
							continue
						} else {
							crl.next = NewRangeLoop(node, tpl, c, w)
							rl = crl.next
							break
						}
					}
				}
				rl.cntr = 0
				rl.node = node
				rl.tpl = tpl
				rl.ctx = c
				rl.w = w
			}
			rl.stat = rlInuse
			c.Err = v.ins.Loop(v.val, rl, &c.buf, c.bufS[1:]...)
			rl.stat = rlFree
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
	c.BufI = cnt
	allowIter = false
	cntr = 0
	for {
		switch node.loopCondOp {
		case OpLt:
			allowIter = c.BufI < lim
		case OpLtq:
			allowIter = c.BufI <= lim
		case OpGt:
			allowIter = c.BufI > lim
		case OpGtq:
			allowIter = c.BufI >= lim
		case OpEq:
			allowIter = c.BufI == lim
		case OpNq:
			allowIter = c.BufI != lim
		default:
			c.Err = ErrWrongLoopCond
			break
		}
		if !allowIter {
			break
		}

		c.Set(fastconv.B2S(node.loopCnt), &c.BufI, &inspector.StaticInspector{})

		if cntr > 0 && len(node.loopSep) > 0 {
			_, _ = w.Write(node.loopSep)
		}
		cntr++
		c.chQB = true
		var err error
		for _, ch := range node.child {
			err = tpl.renderNode(w, ch, c)
			if err == ErrBreakLoop || err == ErrContLoop {
				break
			}
		}
		c.chQB = false

		switch node.loopCntOp {
		case OpInc:
			c.BufI++
		case OpDec:
			c.BufI--
		default:
			c.Err = ErrWrongLoopOp
			break
		}

		if err == ErrBreakLoop {
			break
		}
		if err == ErrContLoop {
			continue
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
		r, ok = if2int(raw)
		if !ok {
			c.Err = ErrWrongLoopLim
			return
		}
	}
	return
}

func (c *Ctx) replaceQB(path []byte) []byte {
	qbLi := bytes.Index(path, qbL)
	qbRi := bytes.Index(path, qbR)
	if qbLi != -1 && qbRi != -1 && qbLi < qbRi && qbRi < len(path) {
		c.Buf.Reset().Write(path[0:qbLi]).Write(dot)
		c.chQB = false
		c.bufX = c.get(path[qbLi+1 : qbRi])
		if c.bufX != nil {
			c.Buf1, c.Err = any2bytes.AnyToBytes(c.Buf1, c.bufX)
			if c.Err != nil {
				c.chQB = true
				return nil
			}
			c.Buf.Write(c.Buf1)
		}
		c.chQB = true
		c.Buf.Write(path[qbRi+1:])
		path = c.Buf
	}
	return path
}
