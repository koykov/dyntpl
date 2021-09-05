package dyntpl

import (
	"bytes"
	"io"
	"strconv"
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/bytebuf"
	"github.com/koykov/fastconv"
	"github.com/koykov/i18n"
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
	// Check json quote/escape/encode flags.
	chJQ, chHE, chUE bool
	// Internal buffers.
	buf   []byte
	bufS  []string
	bufI  int
	bufX  interface{}
	bufA  []interface{}
	bufLC []int64
	bufMO bytebuf.ChainBuf
	// Range loop helper.
	rl *RangeLoop

	// Break depth.
	brkD int

	// List of internal byte writers to process include expressions.
	w  []bytes.Buffer
	wl int

	// List of internal KV pairs.
	kv  []ctxKV
	kvl int

	// i18n support.
	loc  string
	i18n unsafe.Pointer
	repl i18n.PlaceholderReplacer

	// External buffers to use in modifier and condition helpers.
	BufAcc bytebuf.AccumulativeBuf

	Buf, Buf1, Buf2 bytebuf.ChainBuf

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
	// Special case: var is counter.
	cntrF bool
	cntr  int

	ins inspector.Inspector
}

// Context key-value pair.
type ctxKV struct {
	k []byte
	v interface{}
}

var (
	// Byte constants.
	qbL = []byte("[")
	qbR = []byte("]")
	dot = []byte(".")
)

// Make new context object.
func NewCtx() *Ctx {
	ctx := Ctx{
		vars: make([]ctxVar, 0),
		bufS: make([]string, 0),
		Buf:  make(bytebuf.ChainBuf, 0),
		Buf1: make(bytebuf.ChainBuf, 0),
		Buf2: make(bytebuf.ChainBuf, 0),
		buf:  make([]byte, 0),
		bufA: make([]interface{}, 0),
	}
	return &ctx
}

// Set the variable to context.
// Inspector ins should be correspond to variable val.
func (c *Ctx) Set(key string, val interface{}, ins inspector.Inspector) {
	for i := 0; i < c.ln; i++ {
		if c.vars[i].key == key {
			// Update existing variable.
			c.vars[i].val = val
			c.vars[i].ins = ins
			return
		}
	}
	// Add new variable.
	if c.ln < len(c.vars) {
		// Use existing item in variable list..
		c.vars[c.ln].key = key
		c.vars[c.ln].val = val
		c.vars[c.ln].ins = ins
	} else {
		// Extend the variable list with new one.
		c.vars = append(c.vars, ctxVar{
			key: key,
			val: val,
			ins: ins,
		})
	}
	// Increase variables count.
	c.ln++
}

// Set static variable to context.
func (c *Ctx) SetStatic(key string, val interface{}) {
	ins, err := inspector.GetInspector("static")
	if err != nil {
		c.Err = err
		return
	}
	c.Set(key, val, ins)
}

// Set bytes as static variable.
//
// See Ctx.Set().
// This is a special case to improve speed.
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
			ins: ins,
		}
		v.buf = append(v.buf, val...)
		c.vars = append(c.vars, v)
	}
	c.ln++
}

// Set int counter as static variable.
//
// See Ctx.Set().
// This is a special case to support counters.
func (c *Ctx) SetCounter(key string, val int) {
	ins, err := inspector.GetInspector("static")
	if err != nil {
		c.Err = err
		return
	}
	for i := 0; i < c.ln; i++ {
		if c.vars[i].key == key {
			c.vars[i].cntrF = true
			c.vars[i].cntr = val
			c.vars[i].ins = ins
			c.vars[i].val = nil
			c.vars[i].buf = c.vars[i].buf[:0]
			return
		}
	}
	if c.ln < len(c.vars) {
		c.vars[c.ln].key = key
		c.vars[c.ln].cntrF = true
		c.vars[c.ln].cntr = val
		c.vars[c.ln].ins = ins
		c.vars[c.ln].val = nil
		c.vars[c.ln].buf = c.vars[c.ln].buf[:0]
	} else {
		v := ctxVar{
			key:   key,
			cntrF: true,
			cntr:  val,
			ins:   ins,
		}
		c.vars = append(c.vars, v)
	}
	c.ln++
}

// Get arbitrary value from the context by path.
//
// See Ctx.get().
// Path syntax: <ctxVrName>[.<Field>[.<NestedField0>[....<NestedFieldN>]]]
// Examples:
// * user.Bio.Birthday
// * staticVar
func (c *Ctx) Get(path string) interface{} {
	return c.get(fastconv.S2B(path))
}

// Get int counter value.
func (c *Ctx) GetCounter(key string) int {
	rawC := c.Get(key)
	if rawC == nil {
		return 0
	}
	if i, ok := rawC.(*int); ok {
		return *i
	}
	return 0
}

// Config i18n locale and database.
func (c *Ctx) I18n(locale string, db *i18n.DB) {
	c.loc = locale
	if db != nil {
		c.i18n = unsafe.Pointer(db)
	}
}

// Bufferize mod output bytes.
func (c *Ctx) BufModOut(buf *interface{}, p []byte) {
	c.bufMO.Reset().Write(p)
	*buf = &c.bufMO
}

// Bufferize mod output string.
func (c *Ctx) BufModStrOut(buf *interface{}, s string) {
	c.bufMO.Reset().WriteStr(s)
	*buf = &c.bufMO
}

// Reset the context.
//
// Made to use together with pools.
func (c *Ctx) Reset() {
	for i := 0; i < c.ln; i++ {
		c.vars[i].cntrF = false
		c.vars[i].val = nil
		c.vars[i].buf = c.vars[i].buf[:0]
	}
	c.ln = 0

	for i := 0; i < c.wl; i++ {
		c.w[i].Reset()
	}
	c.wl = 0

	c.kvl = 0

	c.loc = ""
	c.i18n = nil
	c.repl.Reset()

	c.Err = nil
	c.bufX = nil
	c.chQB, c.chJQ, c.chHE, c.chUE = false, false, false, false
	c.bufS = c.bufS[:0]
	c.BufAcc.Reset()
	c.bufMO.Reset()
	c.Buf.Reset()
	c.Buf1.Reset()
	c.Buf2.Reset()
	c.buf = c.buf[:0]
	c.bufA = c.bufA[:0]
	c.bufLC = c.bufLC[:0]
	c.brkD = 0
	if c.rl != nil {
		c.rl.Reset()
	}
}

// Internal getter.
//
// See Ctx.Get().
func (c *Ctx) get(path []byte) interface{} {
	// Reset error to avoid catching errors from previous nodes.
	c.Err = nil

	// Special case: check square brackets on counter loops.
	// See Ctx.replaceQB().
	if c.chQB {
		path = c.replaceQB(path)
	}

	// Split path to separate words using dot as separator.
	// So, path user.Bio.Birthday will convert to []string{"user", "Bio", "Birthday"}
	c.bufS = c.bufS[:0]
	c.bufS = bytealg.AppendSplitStr(c.bufS, fastconv.B2S(path), ".", -1)
	if len(c.bufS) == 0 {
		return nil
	}

	// Look for first path chunk in vars.
	for i, v := range c.vars {
		if i == c.ln {
			// Vars limit reached, exit.
			break
		}
		if v.key == c.bufS[0] {
			// Var found.
			if v.val == nil && len(v.buf) > 0 {
				// Special case: var is a byte slice.
				c.Buf.Reset().Write(v.buf)
				c.bufX = &c.Buf
				return c.bufX
			}
			if v.val == nil && v.cntrF {
				// Special case: var is a counter.
				c.bufI = v.cntr
				c.bufX = &c.bufI
				return c.bufX
			}
			// Inspect variable using inspector object.
			// Give search path as list of splitted path minus first key, e.g. []string{"Bio", "Birthday"}
			c.Err = v.ins.GetTo(v.val, &c.bufX, c.bufS[1:]...)
			if c.Err != nil {
				return nil
			}
			return c.bufX
		}
	}

	return nil
}

// Compare method.
func (c *Ctx) cmp(path []byte, cond Op, right []byte) bool {
	// Split path.
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
			// Compare var with right value using inspector.
			if v.cntrF {
				c.Err = v.ins.Cmp(v.cntr, inspector.Op(cond), fastconv.B2S(right), &c.BufB, c.bufS[1:]...)
			} else {
				c.Err = v.ins.Cmp(v.val, inspector.Op(cond), fastconv.B2S(right), &c.BufB, c.bufS[1:]...)
			}
			if c.Err != nil {
				return false
			}
			return c.BufB
		}
	}

	return false
}

// Range loop method to evaluate expressions like:
// {% for k, v := range user.History %}...{% endfor %}
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
			// Look for free range loop object in single-ordered list, see RangeLoop.
			var rl *RangeLoop
			if c.rl == nil {
				// No range loops, create new one.
				c.rl = NewRangeLoop(node, tpl, c, w)
				rl = c.rl
			} else {
				// Move forward over the list while new RL will found.
				crl := c.rl
				for {
					if crl.stat == rlFree {
						// Found it.
						rl = crl
						break
					}
					if crl.stat != rlFree {
						// RL is in use, need to go deeper.
						if crl.next != nil {
							crl = crl.next
							continue
						} else {
							// End of the list, create new free RL and exit from the loop.
							crl.next = NewRangeLoop(node, tpl, c, w)
							rl = crl.next
							break
						}
					}
				}
				// Prepare RL object.
				rl.cntr = 0
				rl.node = node
				rl.tpl = tpl
				rl.ctx = c
				rl.w = w
			}
			// Mark RL as inuse and loop over var using inspector.
			rl.stat = rlInuse
			c.Err = v.ins.Loop(v.val, rl, &c.buf, c.bufS[1:]...)
			rl.stat = rlFree
			return
		}
	}
}

// Counter loop method to evaluate expressions like:
// {% for i:=0; i<10; i++ %}...{% endfor %}
func (c *Ctx) cloop(node Node, tpl *Tpl, w io.Writer) {
	var (
		cnt, lim  int64
		cntr      int
		allowIter bool
	)
	// Prepare bounds.
	cnt = c.cloopRange(node.loopCntStatic, node.loopCntInit)
	if c.Err != nil {
		return
	}
	lim = c.cloopRange(node.loopLimStatic, node.loopLim)
	if c.Err != nil {
		return
	}
	// Prepare counters.
	c.bufLC = append(c.bufLC, cnt)
	idxLC := len(c.bufLC) - 1
	valLC := cnt
	// Start the loop.
	allowIter = false
	cntr = 0
	for {
		// Check iteration allowance.
		switch node.loopCondOp {
		case OpLt:
			allowIter = valLC < lim
		case OpLtq:
			allowIter = valLC <= lim
		case OpGt:
			allowIter = valLC > lim
		case OpGtq:
			allowIter = valLC >= lim
		case OpEq:
			allowIter = valLC == lim
		case OpNq:
			allowIter = valLC != lim
		default:
			c.Err = ErrWrongLoopCond
			break
		}
		// Check breakN signal from child loop.
		allowIter = allowIter && c.brkD == 0

		if !allowIter {
			break
		}

		// Set/update counter var.
		c.SetStatic(fastconv.B2S(node.loopCnt), &c.bufLC[idxLC])

		// Write separator.
		if cntr > 0 && len(node.loopSep) > 0 {
			_, _ = w.Write(node.loopSep)
		}
		cntr++
		// Loop over child nodes with square brackets check in paths.
		c.chQB = true
		var err, lerr error
		for _, ch := range node.child {
			err = tpl.renderNode(w, ch, c)
			if err == ErrLBreakLoop {
				lerr = err
			}
			if err == ErrBreakLoop || err == ErrContLoop {
				break
			}
		}
		c.chQB = false

		// Modify counter var.
		switch node.loopCntOp {
		case OpInc:
			valLC++
			c.bufLC[idxLC]++
		case OpDec:
			valLC--
			c.bufLC[idxLC]--
		default:
			c.Err = ErrWrongLoopOp
			break
		}

		// Handle break/continue cases.
		if err == ErrBreakLoop || lerr == ErrLBreakLoop {
			if c.brkD > 0 {
				c.brkD--
			}
			break
		}
		if err == ErrContLoop {
			continue
		}
	}
	return
}

// Counter loop bound check helper.
//
// Converts initial and final values of the counter to static int values.
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

// Replaces square brackets with variable to concrete values, example:
// user.History[i] -> user.History.0, user.History.1, ...
// , since inspector doesn't supports variadic paths.
func (c *Ctx) replaceQB(path []byte) []byte {
	qbLi := bytes.Index(path, qbL)
	qbRi := bytes.Index(path, qbR)
	if qbLi != -1 && qbRi != -1 && qbLi < qbRi && qbRi < len(path) {
		c.BufAcc.StakeOut()
		c.BufAcc.Write(path[0:qbLi]).Write(dot)
		c.chQB = false
		c.bufX = c.get(path[qbLi+1 : qbRi])
		if c.bufX != nil {
			if err := c.BufAcc.WriteX(c.bufX).Error(); err != nil {
				c.Err = err
				c.chQB = true
				return nil
			}
		}
		c.chQB = true
		c.BufAcc.Write(path[qbRi+1:])
		path = c.BufAcc.StakedBytes()
	}
	return path
}

// Get new or existing byte writer.
//
// Made to write output of including sub-templates.
func (c *Ctx) getW() *bytes.Buffer {
	if c.wl < len(c.w) {
		b := &c.w[c.wl]
		c.wl++
		return b
	} else {
		c.w = append(c.w, bytes.Buffer{})
		b := &c.w[len(c.w)-1]
		c.wl++
		return b
	}
}

// Get new or existing KV pair.
func (c *Ctx) getKV() *ctxKV {
	if c.kvl < len(c.kv) {
		kv := &c.kv[c.kvl]
		c.kvl++
		return kv
	} else {
		c.kv = append(c.kv, ctxKV{})
		kv := &c.kv[len(c.kv)-1]
		c.kvl++
		return kv
	}
}
