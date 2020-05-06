package cbytetpl

import (
	"io"

	"github.com/koykov/fastconv"
	"github.com/koykov/inspector"
)

type RangeLoop struct {
	cntr int
	node Node
	tpl  *Tpl
	ctx  *Ctx
	w    io.Writer
}

func NewRangeLoop(node Node, tpl *Tpl, ctx *Ctx, w io.Writer) *RangeLoop {
	rl := RangeLoop{
		node: node,
		tpl:  tpl,
		ctx:  ctx,
		w:    w,
	}
	return &rl
}

func (rl *RangeLoop) RequireKey() bool {
	return len(rl.node.loopKey) > 0
}

func (rl *RangeLoop) SetKey(val interface{}, ins inspector.Inspector) {
	rl.ctx.Set(fastconv.B2S(rl.node.loopKey), val, ins)
}

func (rl *RangeLoop) SetVal(val interface{}, ins inspector.Inspector) {
	rl.ctx.Set(fastconv.B2S(rl.node.loopVal), val, ins)
}

func (rl *RangeLoop) Loop() {
	if rl.cntr > 0 && len(rl.node.loopSep) > 0 {
		_, _ = rl.w.Write(rl.node.loopSep)
	}
	rl.cntr++
	for _, ch := range rl.node.child {
		_ = rl.tpl.renderNode(rl.w, &ch, rl.ctx)
	}
}
