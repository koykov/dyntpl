package dyntpl

import (
	"bytes"
	"io"
	"sync"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
	"github.com/koykov/x2bytes"
)

// Main template object.
// Template contains only parsed template and evaluation logic.
// All temporary and intermediate data should be store in context object to make using of templates thread-safe.
type Tpl struct {
	Id   string
	tree *Tree
}

var (
	// Templates registry.
	mux         sync.Mutex
	tplRegistry = map[string]*Tpl{}

	// Suppress go vet warning.
	_ = RenderFb
)

// Register template in the registry.
//
// This function can be used in any time to register new templates or overwrite existing to provide dynamics.
func RegisterTpl(id string, tree *Tree) {
	tpl := Tpl{
		Id:   id,
		tree: tree,
	}
	mux.Lock()
	tplRegistry[id] = &tpl
	mux.Unlock()
}

// Render template with id according given context.
//
// See RenderTo().
// Recommend to use RenderTo() together with byte buffer pool to avoid redundant allocations.
func Render(id string, ctx *Ctx) ([]byte, error) {
	buf := bytes.Buffer{}
	err := RenderTo(&buf, id, ctx)
	return buf.Bytes(), err
}

// Render template using fallback id.
//
// See RenderFbTo().
// Using this func you can handle cases when some objects has custom templates and all other should use default templates.
// Example:
// template registry:
// * tplUser
// * tplUser-15
// user object with id 15
// Call of dyntpl.RenderFb("tplUser-15", "tplUser", ctx) will take template tplUser-15 from registry.
// In other case, for user #4:
// call of dyntpl.RenderFbTo("tplUser-4", "tplUser", ctx) will take default template tplUser from registry.
// Recommend to user RenderFbTo().
func RenderFb(id, fbId string, ctx *Ctx) ([]byte, error) {
	buf := bytes.Buffer{}
	err := RenderFbTo(&buf, id, fbId, ctx)
	return buf.Bytes(), err
}

// Render template to given writer object.
//
// Using this function together with byte buffer pool reduces allocations.
func RenderTo(w io.Writer, id string, ctx *Ctx) (err error) {
	mux.Lock()
	tpl, ok := tplRegistry[id]
	mux.Unlock()
	if !ok {
		err = ErrTplNotFound
		return
	}
	return render(w, tpl, ctx)
}

// Render template using fallback ID logic and write result to writer object.
//
// See RenderFb().
// Use this function together with byte buffer pool to reduce allocations.
func RenderFbTo(w io.Writer, id, fbId string, ctx *Ctx) (err error) {
	var (
		tpl *Tpl
		ok  bool
	)
	mux.Lock()
	tpl, ok = tplRegistry[id]
	if !ok {
		tpl, ok = tplRegistry[fbId]
	}
	mux.Unlock()
	if !ok {
		err = ErrTplNotFound
		return
	}
	return render(w, tpl, ctx)
}

// Internal renderer.
func render(w io.Writer, tpl *Tpl, ctx *Ctx) (err error) {
	// Walk over root nodes in tree and evaluate them.
	for _, node := range tpl.tree.nodes {
		err = tpl.renderNode(w, node, ctx)
		if err != nil {
			if err == ErrInterrupt {
				// Interrupt logic.
				err = nil
			}
			return
		}
	}

	return
}

// General node renderer.
func (t *Tpl) renderNode(w io.Writer, node Node, ctx *Ctx) (err error) {
	switch node.typ {
	case TypeRaw:
		if ctx.chJQ {
			// JSON quote mode.
			ctx.Buf.Reset().Write(node.raw)
			ctx.Buf1 = jsonEscape(node.raw, ctx.Buf1)
			_, err = w.Write(ctx.Buf1.Bytes())
		} else if ctx.chHE {
			// HTML escape mode.
			ctx.Buf.Reset().Write(node.raw)
			err = modHtmlEscape(ctx, &ctx.bufX, &ctx.Buf, nil)
			if err != nil {
				_, err = w.Write(node.raw)
			} else {
				_, err = w.Write(ctx.bufX.(*bytealg.ChainBuf).Bytes())
			}
		} else if ctx.chUE {
			// URL encode mode.
			ctx.Buf.Reset().Write(node.raw)
			err = modUrlEncode(ctx, &ctx.bufX, &ctx.Buf, nil)
			if err != nil {
				_, err = w.Write(node.raw)
			} else {
				_, err = w.Write(ctx.bufX.(*bytealg.ChainBuf).Bytes())
			}
		} else {
			// Raw node writes as is.
			_, err = w.Write(node.raw)
		}
	case TypeTpl:
		// Get data from the context.
		raw := ctx.get(node.raw)
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
		// Process modifiers.
		if len(node.mod) > 0 {
			for _, mod := range node.mod {
				// Collect arguments to buffer.
				ctx.bufA = ctx.bufA[:0]
				if len(mod.arg) > 0 {
					for _, arg := range mod.arg {
						if len(arg.name) > 0 {
							kv := ctx.getKV()
							kv.k = arg.name
							if arg.static {
								kv.v = &arg.val
							} else {
								kv.v = ctx.get(arg.val)
							}
							ctx.bufA = append(ctx.bufA, kv)
						} else {
							if arg.static {
								ctx.bufA = append(ctx.bufA, &arg.val)
							} else {
								val := ctx.get(arg.val)
								ctx.bufA = append(ctx.bufA, val)
							}
						}
					}
				}
				ctx.bufX = raw
				// Call the modifier func.
				ctx.Err = (*mod.fn)(ctx, &ctx.bufX, ctx.bufX, ctx.bufA)
				if ctx.Err != nil {
					break
				}
				raw = ctx.bufX
			}
		}
		if ctx.Err != nil {
			return
		}
		if raw == nil || raw == "" {
			// Variable doesn't exists or empty. Do nothing.
			return
		}
		// Convert modified data to bytes array.
		ctx.Buf, err = x2bytes.ToBytesWR(ctx.Buf, raw)
		if err == nil {
			if len(node.prefix) > 0 {
				// Write prefix.
				_, _ = w.Write(node.prefix)
			}
			// Write bytes data.
			_, err = w.Write(ctx.Buf)
			// Write suffix.
			if len(node.suffix) > 0 {
				_, _ = w.Write(node.suffix)
			}
		}
	case TypeCtx:
		// Context node sets new variable, example:
		// {% ctx name = user.Name %} or {% ctx limit = 10 %}

		// It's a speed improvement trick.
		if node.ctxSrcStatic {
			ctx.SetBytes(fastconv.B2S(node.ctxVar), node.ctxSrc)
		} else {
			// Get the inspector.
			ins, err := GetInspector(fastconv.B2S(node.ctxVar), fastconv.B2S(node.ctxIns))
			if err != nil {
				return err
			}

			raw := ctx.get(node.ctxSrc)
			if ctx.Err != nil {
				err = ctx.Err
				return err
			}
			// Process modifiers.
			if len(node.mod) > 0 {
				for _, mod := range node.mod {
					// Collect arguments to buffer.
					ctx.bufA = ctx.bufA[:0]
					if len(mod.arg) > 0 {
						for _, arg := range mod.arg {
							if len(arg.name) > 0 {
								kv := ctx.getKV()
								kv.k = arg.name
								if arg.static {
									kv.v = &arg.val
								} else {
									kv.v = ctx.get(arg.val)
								}
								ctx.bufA = append(ctx.bufA, kv)
							} else {
								if arg.static {
									ctx.bufA = append(ctx.bufA, &arg.val)
								} else {
									val := ctx.get(arg.val)
									ctx.bufA = append(ctx.bufA, val)
								}
							}
						}
					}
					ctx.bufX = raw
					// Call the modifier func.
					ctx.Err = (*mod.fn)(ctx, &ctx.bufX, ctx.bufX, ctx.bufA)
					if ctx.Err != nil {
						break
					}
					raw = ctx.bufX
				}
			}
			if ctx.Err != nil {
				err = ctx.Err
				return err
			}
			empty := raw == nil || raw == ""
			if len(node.ctxOK) > 0 {
				ctx.SetStatic(fastconv.B2S(node.ctxOK), !empty)
			}

			if empty {
				// Empty value, nothing to set. Do nothing and exit.
				return err
			}

			if b, ok := ConvBytes(raw); ok && len(b) > 0 {
				// Set byte array as bytes variable if possible.
				ctx.SetBytes(fastconv.B2S(node.ctxVar), b)
			} else {
				ctx.Set(fastconv.B2S(node.ctxVar), raw, ins)
			}
		}
	case TypeCounter:
		if node.cntrInitF {
			ctx.SetCounter(fastconv.B2S(node.cntrVar), node.cntrInit)
		} else {
			raw := ctx.get(node.cntrVar)
			if ctx.Err != nil {
				err = ctx.Err
				return
			}
			var cntr int
			if cntr64, ok := ConvInt(raw); ok {
				cntr = int(cntr64)
			}
			if node.cntrOp == OpInc {
				cntr += node.cntrOpArg
			} else {
				cntr -= node.cntrOpArg
			}
			ctx.SetCounter(fastconv.B2S(node.cntrVar), cntr)
		}
	case TypeCondOK:
		// Condition-OK node evaluates expressions like if-ok with helper.
		var r bool
		// Check condition-OK helper (mandatory at all).
		if len(node.condHlp) > 0 {
			fn := GetCondOKFn(fastconv.B2S(node.condHlp))
			if fn == nil {
				err = ErrCondHlpNotFound
				return
			}
			// Prepare arguments list.
			ctx.bufA = ctx.bufA[:0]
			if len(node.condHlpArg) > 0 {
				for _, arg := range node.condHlpArg {
					if arg.static {
						ctx.bufA = append(ctx.bufA, &arg.val)
					} else {
						val := ctx.get(arg.val)
						ctx.bufA = append(ctx.bufA, val)
					}
				}
			}
			// Call condition-ok helper func.
			(*fn)(ctx, &ctx.bufX, &ctx.BufB, ctx.bufA)
			r = ctx.BufB
			// Set var, ok to context.
			lv, lr := fastconv.B2S(node.condOKL), fastconv.B2S(node.condOKR)
			ins, err := GetInspector(lv, fastconv.B2S(node.condIns))
			if err != nil {
				return err
			}
			raw := ctx.bufX
			ctx.Set(lv, raw, ins)
			ctx.SetStatic(lr, ctx.BufB)

			// Check extended condition (eg: !ok).
			if len(node.condR) > 0 {
				r, err = t.nodeCmp(&node, ctx)
			}
			// Evaluate condition.
			if r {
				// True case.
				if len(node.child) > 0 {
					err = t.renderNode(w, node.child[0], ctx)
				}
			} else {
				// Else case.
				if len(node.child) > 1 {
					err = t.renderNode(w, node.child[1], ctx)
				}
			}
		}
	case TypeCond:
		// Condition node evaluates condition expressions.
		var r bool
		if len(node.condHlp) > 0 {
			// Condition helper caught.
			fn := GetCondFn(fastconv.B2S(node.condHlp))
			if fn == nil {
				err = ErrCondHlpNotFound
				return
			}
			// Prepare arguments list.
			ctx.bufA = ctx.bufA[:0]
			if len(node.condHlpArg) > 0 {
				for _, arg := range node.condHlpArg {
					if arg.static {
						ctx.bufA = append(ctx.bufA, &arg.val)
					} else {
						val := ctx.get(arg.val)
						ctx.bufA = append(ctx.bufA, val)
					}
				}
			}
			// Call condition helper func.
			r = (*fn)(ctx, ctx.bufA)
		} else {
			r, err = t.nodeCmp(&node, ctx)
		}
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
		// Evaluate condition.
		if r {
			// True case.
			if len(node.child) > 0 {
				err = t.renderNode(w, node.child[0], ctx)
			}
		} else {
			// Else case.
			if len(node.child) > 1 {
				err = t.renderNode(w, node.child[1], ctx)
			}
		}
	case TypeCondTrue, TypeCondFalse, TypeCase, TypeDefault:
		// Just walk over child nodes.
		for _, ch := range node.child {
			err = t.renderNode(w, ch, ctx)
			if err != nil {
				return
			}
		}
	case TypeLoopCount:
		// Evaluate counter loops.
		// See Ctx.cloop().
		ctx.brkD = 0
		ctx.cloop(node, t, w)
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
	case TypeLoopRange:
		// Evaluate range loops.
		// See Ctx.rloop().
		ctx.brkD = 0
		ctx.rloop(node.loopSrc, node, t, w)
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
	case TypeBreak:
		// Break the loop.
		ctx.brkD = node.loopBrkD
		err = ErrBreakLoop
	case TypeLBreak:
		// Lazy break the loop.
		ctx.brkD = node.loopBrkD
		err = ErrLBreakLoop
	case TypeContinue:
		// Go to next iteration of loop.
		err = ErrContLoop
	case TypeSwitch:
		// Switch magic...
		r := false
		if len(node.switchArg) > 0 {
			// Classic switch case.
			for _, ch := range node.child {
				if ch.typ == TypeCase {
					if ch.caseStaticL {
						r = ctx.cmp(node.switchArg, OpEq, ch.caseL)
					} else {
						ctx.get(ch.caseL)
						if ctx.Err == nil {
							ctx.Buf, err = x2bytes.ToBytesWR(ctx.Buf, ctx.bufX)
							if err != nil {
								return
							}
							r = ctx.cmp(node.switchArg, OpEq, ctx.Buf)
						}
					}
				}
				if r {
					err = t.renderNode(w, ch, ctx)
					break
				}
			}
		} else {
			// Switch without condition case.
			for _, ch := range node.child {
				if ch.typ == TypeCase {
					if len(ch.caseHlp) > 0 {
						// Case condition helper caught.
						fn := GetCondFn(fastconv.B2S(ch.caseHlp))
						if fn == nil {
							err = ErrCondHlpNotFound
							return
						}
						// Prepare arguments list.
						ctx.bufA = ctx.bufA[:0]
						if len(ch.caseHlpArg) > 0 {
							for _, arg := range ch.caseHlpArg {
								if arg.static {
									ctx.bufA = append(ctx.bufA, &arg.val)
								} else {
									val := ctx.get(arg.val)
									ctx.bufA = append(ctx.bufA, val)
								}
							}
						}
						// Call condition helper func.
						r = (*fn)(ctx, ctx.bufA)
					} else {
						sl := ch.caseStaticL
						sr := ch.caseStaticR
						if sl && sr {
							err = ErrSenselessCond
							return
						}
						if sr {
							// Right side is static.
							r = ctx.cmp(ch.caseL, ch.caseOp, ch.caseR)
						} else if sl {
							// Left side is static.
							r = ctx.cmp(ch.caseR, ch.caseOp.Swap(), ch.caseL)
						} else {
							// Both sides isn't static.
							ctx.get(ch.caseR)
							if ctx.Err == nil {
								ctx.Buf, err = x2bytes.ToBytesWR(ctx.Buf, ctx.bufX)
								if err != nil {
									return
								}
								r = ctx.cmp(ch.caseL, ch.caseOp, ctx.Buf)
							}
						}
					}
					if ctx.Err != nil {
						err = ctx.Err
						return
					}
					if r {
						err = t.renderNode(w, ch, ctx)
						break
					}
				}
			}
		}
		if !r {
			for _, ch := range node.child {
				if ch.typ == TypeDefault {
					err = t.renderNode(w, ch, ctx)
					break
				}
			}
		}
	case TypeInclude:
		// Include sub-template expression.
		var tpl *Tpl
		mux.Lock()
		for i := 0; i < len(node.tpl); i++ {
			if t, ok := tplRegistry[fastconv.B2S(node.tpl[i])]; ok {
				tpl = t
				break
			}
		}
		mux.Unlock()
		if tpl != nil {
			w1 := ctx.getW()
			if err = render(w1, tpl, ctx); err != nil {
				return
			}

			_, err = w.Write(w1.Bytes())
		} else {
			err = ErrTplNotFound
		}
	case TypeLocale:
		if len(node.loc) > 0 {
			ctx.loc = fastconv.B2S(node.loc)
		}
	case TypeExit:
		// Interrupt template evaluation.
		err = ErrInterrupt
	case TypeJsonQ:
		ctx.chJQ = true
	case TypeEndJsonQ:
		ctx.chJQ = false
	case TypeHtmlE:
		ctx.chHE = true
	case TypeEndHtmlE:
		ctx.chHE = false
	case TypeUrlEnc:
		ctx.chUE = true
	case TypeEndUrlEnc:
		ctx.chUE = false
	default:
		// Unknown node type caught.
		err = ErrUnknownCtl
	}
	return
}

// Evaluate condition expressions.
func (t *Tpl) nodeCmp(node *Node, ctx *Ctx) (r bool, err error) {
	// Regular comparison.
	sl := node.condStaticL
	sr := node.condStaticR
	if sl && sr {
		// It's senseless to compare two static values.
		err = ErrSenselessCond
		return
	}
	if sr {
		// Right side is static. This is a prefer case
		r = ctx.cmp(node.condL, node.condOp, node.condR)
	} else if sl {
		// Left side is static.
		// dyntpl can't handle expressions like {% if 10 > item.Weight %}...
		// therefore it inverts condition to {% if item.Weight < 10 %}...
		r = ctx.cmp(node.condR, node.condOp.Swap(), node.condL)
	} else {
		// Both sides isn't static. This is a bad case, since need to inspect variables twice.
		ctx.get(node.condR)
		if ctx.Err == nil {
			ctx.Buf, err = x2bytes.ToBytesWR(ctx.Buf, ctx.bufX)
			if err != nil {
				return
			}
			r = ctx.cmp(node.condL, node.condOp, ctx.Buf)
		}
	}
	return
}
