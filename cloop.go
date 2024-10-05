package dyntpl

import (
	"io"
	"strconv"

	"github.com/koykov/byteconv"
)

// Counter loop method to evaluate expressions like:
// {% for i:=0; i<10; i++ %}...{% endfor %}
func (ctx *Ctx) cloop(node *node, tpl *Tpl, w io.Writer) {
	var (
		cnt, lim  int64
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
	var c int
	for {
		// Check iteration allowance.
		switch node.loopCondOp {
		case opLt:
			allowIter = valLC < lim
		case opLtq:
			allowIter = valLC <= lim
		case opGt:
			allowIter = valLC > lim
		case opGtq:
			allowIter = valLC >= lim
		case opEq:
			allowIter = valLC == lim
		case opNq:
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
		ctx.SetStatic(byteconv.B2S(node.loopCnt), &ctx.bufLC[idxLC])

		// Write separator.
		if c > 0 && len(node.loopSep) > 0 {
			_, _ = w.Write(node.loopSep)
		}
		c++
		// Loop over child nodes with square brackets check in paths.
		ctx.chQB = true
		var err, lerr error
		child := node.child
		if len(child) > 0 && child[0].typ == typeCondTrue {
			child = child[0].child
		}
		for i := 0; i < len(child); i++ {
			ch := &child[i]
			err = tpl.writeNode(w, ch, ctx)
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
		case opInc:
			valLC++
			ctx.bufLC[idxLC]++
		case opDec:
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

	if c == 0 && len(node.child) > 1 && node.child[1].typ == typeCondFalse {
		child := node.child[1].child
		for j := 0; j < len(child); j++ {
			ch := &child[j]
			if ctx.Err = tpl.writeNode(w, ch, ctx); ctx.Err != nil {
				break
			}
		}
	}

	return
}

// Counter loop bound check helper.
//
// Converts initial and final values of the counter to static int values.
func (ctx *Ctx) cloopRange(static bool, b []byte) (r int64) {
	if static {
		r, ctx.Err = strconv.ParseInt(byteconv.B2S(b), 0, 0)
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
