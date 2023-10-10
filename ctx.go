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

// Ctx is a context object. Contains list of variables available to inspect.
// In addition, has buffers to help develop new helpers without allocations.
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
	bufX  any
	bufA  []any
	bufLC []int64
	bufMO bytebuf.ChainBuf
	bufCB bytebuf.ChainBuf
	// Range loop helper.
	rl *RangeLoop

	// Deferred functions pool.
	dfr []func() error

	// List of variables taken from ipools and registered to return back.
	ipv  []ipoolVar
	ipvl int

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
	// todo remove as unused later
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
	val any
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
	v any
}

var (
	// Byte constants.
	qbL = []byte("[")
	qbR = []byte("]")
	dot = []byte(".")
)

// NewCtx makes new context object.
func NewCtx() *Ctx {
	ctx := Ctx{
		vars: make([]ctxVar, 0),
		bufS: make([]string, 0),
		Buf:  make(bytebuf.ChainBuf, 0),
		Buf1: make(bytebuf.ChainBuf, 0),
		Buf2: make(bytebuf.ChainBuf, 0),
		buf:  make([]byte, 0),
		bufA: make([]any, 0),
	}
	return &ctx
}

// Set the variable to context.
// Inspector ins should be corresponded to variable val.
func (ctx *Ctx) Set(key string, val any, ins inspector.Inspector) {
	for i := 0; i < ctx.ln; i++ {
		if ctx.vars[i].key == key {
			// Update existing variable.
			ctx.vars[i].val = val
			ctx.vars[i].ins = ins
			return
		}
	}
	// Add new variable.
	if ctx.ln < len(ctx.vars) {
		// Use existing item in variable list.
		ctx.vars[ctx.ln].key = key
		ctx.vars[ctx.ln].val = val
		ctx.vars[ctx.ln].ins = ins
	} else {
		// Extend the variable list with new one.
		ctx.vars = append(ctx.vars, ctxVar{
			key: key,
			val: val,
			ins: ins,
		})
	}
	// Increase variables count.
	ctx.ln++
}

// SetStatic sets static variable to context.
func (ctx *Ctx) SetStatic(key string, val any) {
	ins, err := inspector.GetInspector("static")
	if err != nil {
		ctx.Err = err
		return
	}
	ctx.Set(key, val, ins)
}

// SetBytes sets bytes as static variable.
//
// See Ctx.Set().
// This is a special case to improve speed.
func (ctx *Ctx) SetBytes(key string, val []byte) {
	ins, err := inspector.GetInspector("static")
	if err != nil {
		ctx.Err = err
		return
	}
	for i := 0; i < ctx.ln; i++ {
		if ctx.vars[i].key == key {
			ctx.vars[i].buf = append(ctx.vars[i].buf[:0], val...)
			ctx.vars[i].ins = ins
			return
		}
	}
	if ctx.ln < len(ctx.vars) {
		ctx.vars[ctx.ln].key = key
		ctx.vars[ctx.ln].buf = append(ctx.vars[ctx.ln].buf[:0], val...)
		ctx.vars[ctx.ln].ins = ins
	} else {
		v := ctxVar{
			key: key,
			ins: ins,
		}
		v.buf = append(v.buf, val...)
		ctx.vars = append(ctx.vars, v)
	}
	ctx.ln++
}

// SetString sets string as static variable.
func (ctx *Ctx) SetString(key, val string) {
	ctx.SetBytes(key, fastconv.S2B(val))
}

// SetCounter sets int counter as static variable.
//
// See Ctx.Set().
// This is a special case to support counters.
func (ctx *Ctx) SetCounter(key string, val int) {
	ins, err := inspector.GetInspector("static")
	if err != nil {
		ctx.Err = err
		return
	}
	for i := 0; i < ctx.ln; i++ {
		if ctx.vars[i].key == key {
			ctx.vars[i].cntrF = true
			ctx.vars[i].cntr = val
			ctx.vars[i].ins = ins
			ctx.vars[i].val = nil
			ctx.vars[i].buf = ctx.vars[i].buf[:0]
			return
		}
	}
	if ctx.ln < len(ctx.vars) {
		ctx.vars[ctx.ln].key = key
		ctx.vars[ctx.ln].cntrF = true
		ctx.vars[ctx.ln].cntr = val
		ctx.vars[ctx.ln].ins = ins
		ctx.vars[ctx.ln].val = nil
		ctx.vars[ctx.ln].buf = ctx.vars[ctx.ln].buf[:0]
	} else {
		v := ctxVar{
			key:   key,
			cntrF: true,
			cntr:  val,
			ins:   ins,
		}
		ctx.vars = append(ctx.vars, v)
	}
	ctx.ln++
}

// Get arbitrary value from the context by path.
//
// See Ctx.get().
// Path syntax: <ctxVrName>[.<Field>[.<NestedField0>[....<NestedFieldN>]]]
// Examples:
// * user.Bio.Birthday
// * staticVar
func (ctx *Ctx) Get(path string) any {
	return ctx.get(fastconv.S2B(path))
}

// GetCounter gets int counter value.
func (ctx *Ctx) GetCounter(key string) int {
	rawC := ctx.Get(key)
	if rawC == nil {
		return 0
	}
	if i, ok := rawC.(*int); ok {
		return *i
	}
	return 0
}

// I18n sets i18n locale and database.
func (ctx *Ctx) I18n(locale string, db *i18n.DB) {
	ctx.loc = locale
	if db != nil {
		ctx.i18n = unsafe.Pointer(db)
	}
}

// BufModOut buffers mod output bytes.
func (ctx *Ctx) BufModOut(buf *any, p []byte) {
	ctx.bufMO.Reset().Write(p)
	*buf = &ctx.bufMO
}

// BufModStrOut buffers mod output string.
func (ctx *Ctx) BufModStrOut(buf *any, s string) {
	ctx.bufMO.Reset().WriteStr(s)
	*buf = &ctx.bufMO
}

// Defer registers new deferred function.
//
// Function will call after finishing template.
// todo: find a way how to avoid closure allocation.
func (ctx *Ctx) Defer(fn func() error) {
	ctx.dfr = append(ctx.dfr, fn)
}

// AcquireFrom receives new variable from given pool and register it to return batch after finish template processing.
func (ctx *Ctx) AcquireFrom(pool string) (any, error) {
	v, err := ipools_.acquire(pool)
	if err != nil {
		return nil, err
	}
	if ctx.ipvl < len(ctx.ipv) {
		ctx.ipv[ctx.ipvl].key = pool
		ctx.ipv[ctx.ipvl].val = v
	} else {
		ctx.ipv = append(ctx.ipv, ipoolVar{key: pool, val: v})
	}
	ctx.ipvl++
	return v, nil
}

// Reset the context.
//
// Made to use together with pools.
func (ctx *Ctx) Reset() {
	for i := 0; i < ctx.ln; i++ {
		ctx.vars[i].cntrF = false
		ctx.vars[i].val = nil
		ctx.vars[i].buf = ctx.vars[i].buf[:0]
	}
	ctx.ln = 0

	for i := 0; i < ctx.wl; i++ {
		ctx.w[i].Reset()
	}
	ctx.wl = 0

	ctx.kvl = 0

	ctx.loc = ""
	ctx.i18n = nil
	ctx.repl.Reset()

	ctx.Err = nil
	ctx.bufX = nil
	ctx.chQB, ctx.chJQ, ctx.chHE, ctx.chUE = false, false, false, false
	ctx.bufS = ctx.bufS[:0]
	ctx.bufCB.Reset()
	ctx.BufAcc.Reset()
	ctx.bufMO.Reset()
	ctx.Buf.Reset()
	ctx.Buf1.Reset()
	ctx.Buf2.Reset()
	ctx.buf = ctx.buf[:0]
	ctx.bufA = ctx.bufA[:0]
	ctx.bufLC = ctx.bufLC[:0]
	ctx.brkD = 0
	if ctx.rl != nil {
		ctx.rl.Reset()
	}

	ctx.dfr = ctx.dfr[:0]

	for i := 0; i < ctx.ipvl; i++ {
		_ = ipools_.release(ctx.ipv[i].key, ctx.ipv[i].val)
		ctx.ipv[i].key, ctx.ipv[i].val = "", nil
	}
	ctx.ipvl = 0
}

// Internal getter.
//
// See Ctx.Get().
func (ctx *Ctx) get(path []byte) any {
	// Reset error to avoid catching errors from previous nodes.
	ctx.Err = nil

	// Special case: check square brackets on counter loops.
	// See Ctx.replaceQB().
	if ctx.chQB {
		path = ctx.replaceQB(path)
	}

	// Split path to separate words using dot as separator.
	// So, path user.Bio.Birthday will convert to []string{"user", "Bio", "Birthday"}
	ctx.bufS = ctx.bufS[:0]
	ctx.bufS = bytealg.AppendSplit(ctx.bufS, fastconv.B2S(path), ".", -1)
	if len(ctx.bufS) == 0 {
		return nil
	}

	// Look for first path chunk in vars.
	for i, v := range ctx.vars {
		if i == ctx.ln {
			// Vars limit reached, exit.
			break
		}
		if v.key == ctx.bufS[0] {
			// Var found.
			if v.val == nil && len(v.buf) > 0 {
				// Special case: var is a byte slice.
				ctx.bufCB.Reset().Write(v.buf)
				ctx.bufX = &ctx.bufCB
				return ctx.bufX
			}
			if v.val == nil && v.cntrF {
				// Special case: var is a counter.
				ctx.bufI = v.cntr
				ctx.bufX = &ctx.bufI
				return ctx.bufX
			}
			// Inspect variable using inspector object.
			// Give search path as list of split path minus first key, e.g. []string{"Bio", "Birthday"}
			ctx.bufX = nil
			ctx.Err = v.ins.GetTo(v.val, &ctx.bufX, ctx.bufS[1:]...)
			if ctx.Err != nil {
				return nil
			}
			return ctx.bufX
		}
	}

	return nil
}

// Compare method.
func (ctx *Ctx) cmp(path []byte, cond Op, right []byte) bool {
	// Split path.
	ctx.bufS = ctx.bufS[:0]
	ctx.bufS = bytealg.AppendSplit(ctx.bufS, fastconv.B2S(path), ".", -1)
	if len(ctx.bufS) == 0 {
		return false
	}

	for i, v := range ctx.vars {
		if i == ctx.ln {
			break
		}
		if v.key == ctx.bufS[0] {
			// Compare var with right value using inspector.
			if v.cntrF {
				ctx.Err = v.ins.Compare(v.cntr, inspector.Op(cond), fastconv.B2S(right), &ctx.BufB, ctx.bufS[1:]...)
			} else {
				ctx.Err = v.ins.Compare(v.val, inspector.Op(cond), fastconv.B2S(right), &ctx.BufB, ctx.bufS[1:]...)
			}
			if ctx.Err != nil {
				return false
			}
			return ctx.BufB
		}
	}

	return false
}

func (ctx *Ctx) cmpLC(lc lc, path []byte, cond Op, right []byte) bool {
	ctx.Err = nil
	if ctx.chQB {
		path = ctx.replaceQB(path)
	}

	ctx.bufS = ctx.bufS[:0]
	ctx.bufS = bytealg.AppendSplit(ctx.bufS, fastconv.B2S(path), ".", -1)
	if len(ctx.bufS) == 0 {
		return false
	}

	for i, v := range ctx.vars {
		if i == ctx.ln {
			break
		}
		if v.key == ctx.bufS[0] {
			switch lc {
			case lcLen:
				ctx.Err = v.ins.Length(v.val, &ctx.bufI, ctx.bufS[1:]...)
			case lcCap:
				ctx.Err = v.ins.Capacity(v.val, &ctx.bufI, ctx.bufS[1:]...)
			default:
				return false
			}
			if ctx.Err != nil {
				return false
			}
			si := inspector.StaticInspector{}
			ctx.BufB = false
			ctx.Err = si.Compare(ctx.bufI, inspector.Op(cond), fastconv.B2S(right), &ctx.BufB)
			return ctx.BufB
		}
	}
	return false
}

// Range loop method to evaluate expressions like:
// {% for k, v := range user.History %}...{% endfor %}
func (ctx *Ctx) rloop(path []byte, node Node, tpl *Tpl, w io.Writer) {
	ctx.bufS = ctx.bufS[:0]
	ctx.bufS = bytealg.AppendSplit(ctx.bufS, fastconv.B2S(path), ".", -1)
	if len(ctx.bufS) == 0 {
		return
	}
	for i, v := range ctx.vars {
		if i == ctx.ln {
			break
		}
		if v.key == ctx.bufS[0] {
			// Look for free-range loop object in single-ordered list, see RangeLoop.
			var rl *RangeLoop
			if ctx.rl == nil {
				// No range loops, create new one.
				ctx.rl = NewRangeLoop(node, tpl, ctx, w)
				rl = ctx.rl
			} else {
				// Move forward over the list while new RL will found.
				crl := ctx.rl
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
							crl.next = NewRangeLoop(node, tpl, ctx, w)
							rl = crl.next
							break
						}
					}
				}
				// Prepare RL object.
				rl.cntr = 0
				rl.node = node
				rl.tpl = tpl
				rl.ctx = ctx
				rl.w = w
			}
			// Mark RL as inuse and loop over var using inspector.
			rl.stat = rlInuse
			ctx.Err = v.ins.Loop(v.val, rl, &ctx.buf, ctx.bufS[1:]...)
			rl.stat = rlFree
			return
		}
	}
}

// Counter loop method to evaluate expressions like:
// {% for i:=0; i<10; i++ %}...{% endfor %}
func (ctx *Ctx) cloop(node Node, tpl *Tpl, w io.Writer) {
	var (
		cnt, lim  int64
		cntr      int
		allowIter bool
	)
	// Prepare bounds.
	cnt = ctx.cloopRange(node.loopCntStatic, node.loopCntInit)
	if ctx.Err != nil {
		return
	}
	lim = ctx.cloopRange(node.loopLimStatic, node.loopLim)
	if ctx.Err != nil {
		return
	}
	// Prepare counters.
	ctx.bufLC = append(ctx.bufLC, cnt)
	idxLC := len(ctx.bufLC) - 1
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
			ctx.Err = ErrWrongLoopCond
			break
		}
		// Check breakN signal from child loop.
		allowIter = allowIter && ctx.brkD == 0

		if !allowIter {
			break
		}

		// Set/update counter var.
		ctx.SetStatic(fastconv.B2S(node.loopCnt), &ctx.bufLC[idxLC])

		// Write separator.
		if cntr > 0 && len(node.loopSep) > 0 {
			_, _ = w.Write(node.loopSep)
		}
		cntr++
		// Loop over child nodes with square brackets check in paths.
		ctx.chQB = true
		var err, lerr error
		for _, ch := range node.child {
			err = tpl.renderNode(w, ch, ctx)
			if err == ErrLBreakLoop {
				lerr = err
			}
			if err == ErrBreakLoop || err == ErrContLoop {
				break
			}
		}
		ctx.chQB = false

		// Modify counter var.
		switch node.loopCntOp {
		case OpInc:
			valLC++
			ctx.bufLC[idxLC]++
		case OpDec:
			valLC--
			ctx.bufLC[idxLC]--
		default:
			ctx.Err = ErrWrongLoopOp
			break
		}

		// Handle break/continue cases.
		if err == ErrBreakLoop || lerr == ErrLBreakLoop {
			if ctx.brkD > 0 {
				ctx.brkD--
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
func (ctx *Ctx) cloopRange(static bool, b []byte) (r int64) {
	if static {
		r, ctx.Err = strconv.ParseInt(fastconv.B2S(b), 0, 0)
		if ctx.Err != nil {
			return
		}
	} else {
		var ok bool
		raw := ctx.get(b)
		if ctx.Err != nil {
			return
		}
		r, ok = if2int(raw)
		if !ok {
			ctx.Err = ErrWrongLoopLim
			return
		}
	}
	return
}

// Replaces square brackets with variable to concrete values, example:
// user.History[i] -> user.History.0, user.History.1, ...
// , since inspector doesn't support variadic paths.
func (ctx *Ctx) replaceQB(path []byte) []byte {
	qbLi := bytes.Index(path, qbL)
	qbRi := bytes.Index(path, qbR)
	if qbLi != -1 && qbRi != -1 && qbLi < qbRi && qbRi < len(path) {
		ctx.BufAcc.StakeOut()
		ctx.BufAcc.Write(path[0:qbLi]).Write(dot)
		ctx.chQB = false
		ctx.bufX = ctx.get(path[qbLi+1 : qbRi])
		if ctx.bufX != nil {
			if err := ctx.BufAcc.WriteX(ctx.bufX).Error(); err != nil {
				ctx.Err = err
				ctx.chQB = true
				return nil
			}
		}
		ctx.chQB = true
		ctx.BufAcc.Write(path[qbRi+1:])
		path = ctx.BufAcc.StakedBytes()
	}
	return path
}

// Get new or existing byte writer.
//
// Made to write output of including sub-templates.
func (ctx *Ctx) getW() *bytes.Buffer {
	var b *bytes.Buffer
	if ctx.wl < len(ctx.w) {
		b = &ctx.w[ctx.wl]
		ctx.wl++
	} else {
		ctx.w = append(ctx.w, bytes.Buffer{})
		b = &ctx.w[len(ctx.w)-1]
		ctx.wl++
	}
	return b
}

// Get new or existing KV pair.
func (ctx *Ctx) getKV() *ctxKV {
	var kv *ctxKV
	if ctx.kvl < len(ctx.kv) {
		kv = &ctx.kv[ctx.kvl]
		ctx.kvl++
	} else {
		ctx.kv = append(ctx.kv, ctxKV{})
		kv = &ctx.kv[len(ctx.kv)-1]
		ctx.kvl++
	}
	return kv
}

func (ctx *Ctx) defer_() (err error) {
	if len(ctx.dfr) > 0 {
		for i := 0; i < len(ctx.dfr); i++ {
			if err = ctx.dfr[i](); err != nil {
				break
			}
		}
	}
	return
}
