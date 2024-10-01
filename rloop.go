package dyntpl

import (
	"io"

	"github.com/koykov/byteconv"
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
	node *Node
	tpl  *Tpl
	ctx  *Ctx
	next *RangeLoop
	c    uint
	w    io.Writer
}

// NewRangeLoop makes new RL.
func NewRangeLoop(node *Node, tpl *Tpl, ctx *Ctx, w io.Writer) *RangeLoop {
	rl := RangeLoop{
		node: node,
		tpl:  tpl,
		ctx:  ctx,
		w:    w,
	}
	return &rl
}

// RequireKey checks if node requires a key to store in the context.
func (rl *RangeLoop) RequireKey() bool {
	return len(rl.node.loopKey) > 0
}

// SetKey saves key to the context.
func (rl *RangeLoop) SetKey(val any, ins inspector.Inspector) {
	rl.ctx.Set(byteconv.B2S(rl.node.loopKey), val, ins)
}

// SetVal saves value to the context.
func (rl *RangeLoop) SetVal(val any, ins inspector.Inspector) {
	rl.ctx.Set(byteconv.B2S(rl.node.loopVal), val, ins)
}

// Iterate performs the iteration.
func (rl *RangeLoop) Iterate() inspector.LoopCtl {
	rl.c++
	if rl.ctx.brkD > 0 {
		return inspector.LoopCtlBrk
	}

	if rl.cntr > 0 && len(rl.node.loopSep) > 0 {
		_, _ = rl.w.Write(rl.node.loopSep)
	}
	rl.cntr++
	var err, lerr error
	child := rl.node.child
	if len(child) > 0 && child[0].typ == typeCondTrue {
		child = child[0].child
	}
	for i := 0; i < len(child); i++ {
		ch := &child[i]
		err = rl.tpl.writeNode(rl.w, ch, rl.ctx)
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

// Reset clears all data in the list of RL.
func (rl *RangeLoop) Reset() {
	crl := rl
	for crl != nil {
		crl.stat = rlFree
		crl.cntr = 0
		crl.ctx = nil
		crl.tpl = nil
		crl.c = 0
		crl.w = nil
		crl = crl.next
	}
}
