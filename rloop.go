package dyntpl

import (
	"io"

	"github.com/koykov/fastconv"
	"github.com/koykov/inspector"
)

const (
	// RangeLoop object statuses.
	rlFree  = uint(0)
	rlInuse = uint(1)
)

// RangeLoop is a object that injects to inspector to perform range loop execution.
type RangeLoop struct {
	cntr int
	stat uint
	node Node
	tpl  *Tpl
	ctx  *Ctx
	next *RangeLoop
	w    io.Writer
}

// Init new RL.
func NewRangeLoop(node Node, tpl *Tpl, ctx *Ctx, w io.Writer) *RangeLoop {
	rl := RangeLoop{
		node: node,
		tpl:  tpl,
		ctx:  ctx,
		w:    w,
	}
	return &rl
}

// Check if node requires a key to store in the context.
func (rl *RangeLoop) RequireKey() bool {
	return len(rl.node.loopKey) > 0
}

// Save key to the context.
func (rl *RangeLoop) SetKey(val interface{}, ins inspector.Inspector) {
	rl.ctx.Set(fastconv.B2S(rl.node.loopKey), val, ins)
}

// Save value to the context.
func (rl *RangeLoop) SetVal(val interface{}, ins inspector.Inspector) {
	rl.ctx.Set(fastconv.B2S(rl.node.loopVal), val, ins)
}

// Perform the iteration.
func (rl *RangeLoop) Iterate() inspector.LoopCtl {
	if rl.ctx.brkD > 0 {
		return inspector.LoopCtlBrk
	}

	if rl.cntr > 0 && len(rl.node.loopSep) > 0 {
		_, _ = rl.w.Write(rl.node.loopSep)
	}
	rl.cntr++
	var err, lerr error
	for _, ch := range rl.node.child {
		err = rl.tpl.renderNode(rl.w, ch, rl.ctx)
		if err == ErrLBreakLoop {
			lerr = err
		}
		if err == ErrBreakLoop {
			if rl.ctx.brkD > 0 {
				rl.ctx.brkD--
			}
			return inspector.LoopCtlBrk
		}
		if err == ErrContLoop {
			return inspector.LoopCtlCnt
		}
	}
	if err == ErrBreakLoop || lerr == ErrLBreakLoop {
		if rl.ctx.brkD > 0 {
			rl.ctx.brkD--
		}
		return inspector.LoopCtlBrk
	}
	return inspector.LoopCtlNone
}

// Clear all data in the list of RL.
func (rl *RangeLoop) Reset() {
	crl := rl
	for crl != nil {
		crl.stat = rlFree
		crl.cntr = 0
		crl.ctx = nil
		crl.tpl = nil
		crl.w = nil
		crl = crl.next
	}
}
