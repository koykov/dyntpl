package dyntpl

import (
	"bytes"
	"io"

	"github.com/koykov/fastconv"
)

// Tpl is a main template object.
// Template contains only parsed template and evaluation logic.
// All temporary and intermediate data should be store in context object to make using of templates thread-safe.
type Tpl struct {
	ID   int
	Key  string
	tree *Tree
}

var (
	// Templates DB.
	tplDB = initDB()

	// Suppress go vet warning.
	_, _, _ = RegisterTplID, RenderFallback, RenderByID
)

// RegisterTpl saves template by ID and key in the registry.
//
// You may use to access to the template both ID or key.
// This function can be used in any time to register new templates or overwrite existing to provide dynamics.
func RegisterTpl(id int, key string, tree *Tree) {
	tplDB.set(id, key, tree)
}

// RegisterTplID saves template using only ID.
//
// See RegisterTpl().
func RegisterTplID(id int, tree *Tree) {
	tplDB.set(id, "-1", tree)
}

// RegisterTplKey saves template using only key.
//
// See RegisterTpl().
func RegisterTplKey(key string, tree *Tree) {
	tplDB.set(-1, key, tree)
}

// Render template with given key according given context.
//
// See Write().
// Recommend to use Write() together with byte buffer pool to avoid redundant allocations.
func Render(key string, ctx *Ctx) ([]byte, error) {
	buf := bytes.Buffer{}
	err := Write(&buf, key, ctx)
	return buf.Bytes(), err
}

// RenderFallback renders template using one of keys: key or fallback key.
//
// See WriteFallback().
// Using this func you can handle cases when some objects have custom templates and all other should use default templates.
// Example:
// template registry:
// * tplUser
// * tplUser-15
// user object with id 15
// Call of dyntpl.RenderFallback("tplUser-15", "tplUser", ctx) will take template tplUser-15 from registry.
// In other case, for user #4:
// call of dyntpl.WriteFallback("tplUser-4", "tplUser", ctx) will take default template tplUser from registry.
// Recommend to user WriteFallback().
func RenderFallback(key, fbKey string, ctx *Ctx) ([]byte, error) {
	buf := bytes.Buffer{}
	err := WriteFallback(&buf, key, fbKey, ctx)
	return buf.Bytes(), err
}

// Write template with given key to given writer object.
//
// Using this function together with byte buffer pool reduces allocations.
func Write(w io.Writer, key string, ctx *Ctx) (err error) {
	tpl := tplDB.getKey(key)
	if tpl == nil {
		err = ErrTplNotFound
		return
	}
	return write(w, tpl, ctx)
}

// WriteFallback writes template using fallback key logic and write result to writer object.
//
// See RenderFallback().
// Use this function together with byte buffer pool to reduce allocations.
func WriteFallback(w io.Writer, key, fbKey string, ctx *Ctx) (err error) {
	tpl := tplDB.getKey1(key, fbKey)
	if tpl == nil {
		err = ErrTplNotFound
		return
	}
	return write(w, tpl, ctx)
}

// RenderByID renders template with given ID according context.
//
// See WriteByID().
// Recommend to use WriteByID() together with byte buffer pool to avoid redundant allocations.
func RenderByID(id int, ctx *Ctx) ([]byte, error) {
	buf := bytes.Buffer{}
	err := WriteByID(&buf, id, ctx)
	return buf.Bytes(), err
}

// WriteByID writes template with given ID to given writer object.
//
// Using this function together with byte buffer pool reduces allocations.
func WriteByID(w io.Writer, id int, ctx *Ctx) (err error) {
	tpl := tplDB.getID(id)
	if tpl == nil {
		err = ErrTplNotFound
		return
	}
	return write(w, tpl, ctx)
}

// Internal renderer.
func write(w io.Writer, tpl *Tpl, ctx *Ctx) (err error) {
	// Walk over root nodes in tree and evaluate them.
	for _, node := range tpl.tree.nodes {
		err = tpl.writeNode(w, node, ctx)
		if err != nil {
			if err == ErrInterrupt {
				// Interrupt logic.
				err = nil
			}
			return
		}
	}

	// Call defer functions consecutively.
	// First failed function will stop that process and return error encountered.
	err = ctx.defer_()

	return
}

// General node renderer.
func (t *Tpl) writeNode(w io.Writer, node Node, ctx *Ctx) (err error) {
	switch node.typ {
	case TypeRaw:
		if ctx.chJQ {
			// JSON quote mode.
			ctx.BufAcc.StakeOut()
			jsonEscape(node.raw, &ctx.BufAcc)
			_, err = w.Write(ctx.BufAcc.StakedBytes())
		} else if ctx.chHE {
			// HTML escape mode.
			ctx.bufCB.Reset().Write(node.raw)
			err = modHTMLEscape(ctx, &ctx.bufX, &ctx.bufCB, nil)
			if err != nil {
				_, err = w.Write(node.raw)
			} else {
				_, err = w.Write(ctx.bufMO.Bytes())
			}
		} else if ctx.chUE {
			// URL encode mode.
			ctx.bufCB.Reset().Write(node.raw)
			err = modURLEncode(ctx, &ctx.bufX, &ctx.bufCB, nil)
			if err != nil {
				_, err = w.Write(node.raw)
			} else {
				_, err = w.Write(ctx.bufMO.Bytes())
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
			for _, mod_ := range node.mod {
				// Collect arguments to buffer.
				ctx.bufA = ctx.bufA[:0]
				if len(mod_.arg) > 0 {
					for _, arg_ := range mod_.arg {
						if len(arg_.name) > 0 {
							kv := ctx.getKV()
							kv.k = arg_.name
							if arg_.static {
								kv.v = &arg_.val
							} else {
								kv.v = ctx.get(arg_.val)
							}
							ctx.bufA = append(ctx.bufA, kv)
						} else {
							if arg_.static {
								ctx.bufA = append(ctx.bufA, &arg_.val)
							} else {
								val := ctx.get(arg_.val)
								ctx.bufA = append(ctx.bufA, val)
							}
						}
					}
				}
				ctx.bufX = raw
				// Call the modifier func.
				ctx.Err = mod_.fn(ctx, &ctx.bufX, ctx.bufX, ctx.bufA)
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
			// Variable doesn't exist or empty. Do nothing.
			return
		}
		// Convert modified data to bytes array.
		if err = ctx.BufAcc.StakeOut().WriteX(raw).Error(); err == nil {
			if len(node.prefix) > 0 {
				// Write prefix.
				_, _ = w.Write(node.prefix)
			}
			// Write bytes data.
			_, err = w.Write(ctx.BufAcc.StakedBytes())
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
				for _, mod_ := range node.mod {
					// Collect arguments to buffer.
					ctx.bufA = ctx.bufA[:0]
					if len(mod_.arg) > 0 {
						for _, arg_ := range mod_.arg {
							if len(arg_.name) > 0 {
								kv := ctx.getKV()
								kv.k = arg_.name
								if arg_.static {
									kv.v = &arg_.val
								} else {
									kv.v = ctx.get(arg_.val)
								}
								ctx.bufA = append(ctx.bufA, kv)
							} else {
								if arg_.static {
									ctx.bufA = append(ctx.bufA, &arg_.val)
								} else {
									val := ctx.get(arg_.val)
									ctx.bufA = append(ctx.bufA, val)
								}
							}
						}
					}
					ctx.bufX = raw
					// Call the modifier func.
					ctx.Err = mod_.fn(ctx, &ctx.bufX, ctx.bufX, ctx.bufA)
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
				for _, arg_ := range node.condHlpArg {
					if arg_.static {
						ctx.bufA = append(ctx.bufA, &arg_.val)
					} else {
						val := ctx.get(arg_.val)
						ctx.bufA = append(ctx.bufA, val)
					}
				}
			}
			// Call condition-ok helper func.
			fn(ctx, &ctx.bufX, &ctx.BufB, ctx.bufA)
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
					err = t.writeNode(w, node.child[0], ctx)
				}
			} else {
				// Else case.
				if len(node.child) > 1 {
					err = t.writeNode(w, node.child[1], ctx)
				}
			}
		}
	case TypeCond:
		// Condition node evaluates condition expressions.
		var r bool
		switch {
		case len(node.condHlp) > 0 && node.condLC == lcNone:
			// Condition helper caught (no LC case).
			fn := GetCondFn(fastconv.B2S(node.condHlp))
			if fn == nil {
				err = ErrCondHlpNotFound
				return
			}
			// Prepare arguments list.
			ctx.bufA = ctx.bufA[:0]
			if len(node.condHlpArg) > 0 {
				for _, arg_ := range node.condHlpArg {
					if arg_.static {
						ctx.bufA = append(ctx.bufA, &arg_.val)
					} else {
						val := ctx.get(arg_.val)
						ctx.bufA = append(ctx.bufA, val)
					}
				}
			}
			// Call condition helper func.
			r = fn(ctx, ctx.bufA)
		case len(node.condHlp) > 0 && node.condLC > lcNone:
			// Condition helper in LC mode.
			if len(node.condHlpArg) == 0 {
				err = ErrModNoArgs
				return
			}
			r = ctx.cmpLC(node.condLC, node.condHlpArg[0].val, node.condOp, node.condR)
		default:
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
				err = t.writeNode(w, node.child[0], ctx)
			}
		} else {
			// Else case.
			if len(node.child) > 1 {
				err = t.writeNode(w, node.child[1], ctx)
			}
		}
	case TypeCondTrue, TypeCondFalse, TypeCase, TypeDefault:
		// Just walk over child nodes.
		for _, ch := range node.child {
			err = t.writeNode(w, ch, ctx)
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
							if err = ctx.BufAcc.StakeOut().WriteX(ctx.bufX).Error(); err != nil {
								return
							}
							r = ctx.cmp(node.switchArg, OpEq, ctx.BufAcc.StakedBytes())
						}
					}
				}
				if r {
					err = t.writeNode(w, ch, ctx)
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
							for _, arg_ := range ch.caseHlpArg {
								if arg_.static {
									ctx.bufA = append(ctx.bufA, &arg_.val)
								} else {
									val := ctx.get(arg_.val)
									ctx.bufA = append(ctx.bufA, val)
								}
							}
						}
						// Call condition helper func.
						r = fn(ctx, ctx.bufA)
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
							// Both sides aren't static.
							ctx.get(ch.caseR)
							if ctx.Err == nil {
								if err = ctx.BufAcc.StakeOut().WriteX(ctx.bufX).Error(); err != nil {
									return
								}
								r = ctx.cmp(ch.caseL, ch.caseOp, ctx.BufAcc.StakedBytes())
							}
						}
					}
					if ctx.Err != nil {
						err = ctx.Err
						return
					}
					if r {
						err = t.writeNode(w, ch, ctx)
						break
					}
				}
			}
		}
		if !r {
			for _, ch := range node.child {
				if ch.typ == TypeDefault {
					err = t.writeNode(w, ch, ctx)
					break
				}
			}
		}
	case TypeInclude:
		// Include sub-template expression.
		tpl := tplDB.getBKeys(node.tpl)
		if tpl != nil {
			w1 := ctx.getW()
			if err = write(w1, tpl, ctx); err != nil {
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
		// Right side is static. This is preferred case
		r = ctx.cmp(node.condL, node.condOp, node.condR)
	} else if sl {
		// Left side is static.
		// dyntpl can't handle expressions like {% if 10 > item.Weight %}...
		// therefore it inverts condition to {% if item.Weight < 10 %}...
		r = ctx.cmp(node.condR, node.condOp.Swap(), node.condL)
	} else {
		// Both sides aren't static. This is a bad case, since need to inspect variables twice.
		ctx.get(node.condR)
		if ctx.Err == nil {
			if err = ctx.BufAcc.StakeOut().WriteX(ctx.bufX).Error(); err != nil {
				return
			}
			r = ctx.cmp(node.condL, node.condOp, ctx.BufAcc.StakedBytes())
		}
	}
	return
}
