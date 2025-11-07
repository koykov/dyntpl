package dyntpl

import (
	"bytes"
	"io"

	"github.com/koykov/byteconv"
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
	for i := 0; i < len(tpl.tree.nodes); i++ {
		n := &tpl.tree.nodes[i]
		err = tpl.writeNode(w, n, ctx)
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
func (t *Tpl) writeNode(w io.Writer, node *node, ctx *Ctx) (err error) {
	switch node.typ {
	case typeRaw:
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
	case typeTpl:
		// Get data from the context.
		raw := ctx.get(node.raw)
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
		ctx.noesc = node.noesc
		// Process modifiers.
		if n := len(node.mod); n > 0 {
			_ = node.mod[n-1]
			for i := 0; i < n; i++ {
				mod_ := &node.mod[i]
				// Collect arguments to buffer.
				ctx.bufA = ctx.bufA[:0]
				if m := len(mod_.arg); m > 0 {
					_ = mod_.arg[m-1]
					for j := 0; j < m; j++ {
						arg_ := mod_.arg[j]
						if len(arg_.name) > 0 {
							kv := ctx.getKV()
							kv.K = arg_.name
							if arg_.static {
								kv.V = &arg_.val
							} else {
								kv.V = ctx.get(arg_.val)
							}
							ctx.bufA = append(ctx.bufA, kv)
						} else {
							if arg_.global {
								ctx.bufA = append(ctx.bufA, GetGlobal(byteconv.B2S(arg_.val)))
							} else if arg_.static {
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
		ctx.noesc = false
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
	case typeCtx:
		// Context node sets new variable, example:
		// {% ctx name = user.Name %} or {% ctx limit = 10 %}

		// It's a speed improvement trick.
		if node.ctxSrcStatic {
			ctx.SetBytes(byteconv.B2S(node.ctxVar), node.ctxSrc)
		} else {
			// Get the inspector.
			ins, err := GetInspector(byteconv.B2S(node.ctxVar), byteconv.B2S(node.ctxIns))
			if err != nil {
				return err
			}

			raw := ctx.get(node.ctxSrc)
			if ctx.Err != nil {
				err = ctx.Err
				return err
			}
			// Process modifiers.
			if n := len(node.mod); n > 0 {
				_ = node.mod[n-1]
				for i := 0; i < n; i++ {
					mod_ := &node.mod[i]
					// Collect arguments to buffer.
					ctx.bufA = ctx.bufA[:0]
					if m := len(mod_.arg); m > 0 {
						_ = mod_.arg[m-1]
						for j := 0; j < len(mod_.arg); j++ {
							arg_ := mod_.arg[j]
							if len(arg_.name) > 0 {
								kv := ctx.getKV()
								kv.K = arg_.name
								if arg_.static {
									kv.V = &arg_.val
								} else {
									kv.V = ctx.get(arg_.val)
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
				ctx.SetStatic(byteconv.B2S(node.ctxOK), !empty)
			}

			if empty {
				// Empty value, nothing to set. Do nothing and exit.
				return err
			}

			if b, ok := ConvBytes(raw); ok && len(b) > 0 {
				// Set byte array as bytes variable if possible.
				ctx.SetBytes(byteconv.B2S(node.ctxVar), b)
			} else {
				ctx.Set(byteconv.B2S(node.ctxVar), raw, ins)
			}
		}
	case typeCounter:
		if node.cntrInitF {
			ctx.SetCounter(byteconv.B2S(node.cntrVar), node.cntrInit)
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
			if node.cntrOp == opInc {
				cntr += node.cntrOpArg
			} else {
				cntr -= node.cntrOpArg
			}
			ctx.SetCounter(byteconv.B2S(node.cntrVar), cntr)
		}
	case typeCondOK:
		// Condition-OK node evaluates expressions like if-ok with helper.
		var r bool
		// Check condition-OK helper (mandatory at all).
		if len(node.condHlp) > 0 {
			fn := GetCondOKFn(byteconv.B2S(node.condHlp))
			if fn == nil {
				err = ErrCondHlpNotFound
				return
			}
			// Prepare arguments list.
			ctx.bufA = ctx.bufA[:0]
			if n := len(node.condHlpArg); n > 0 {
				_ = node.condHlpArg[n-1]
				for i := 0; i < n; i++ {
					arg_ := node.condHlpArg[i]
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
			lv, lr := byteconv.B2S(node.condOKL), byteconv.B2S(node.condOKR)
			ins, err := GetInspector(lv, byteconv.B2S(node.condIns))
			if err != nil {
				return err
			}
			raw := ctx.bufX
			ctx.Set(lv, raw, ins)
			ctx.SetStatic(lr, ctx.BufB)

			// Check extended condition (eg: !ok).
			if len(node.condR) > 0 {
				r, err = t.nodeCmp(node, ctx)
			}
			// Evaluate condition.
			if r {
				// True case.
				if len(node.child) > 0 {
					err = t.writeNode(w, &node.child[0], ctx)
				}
			} else {
				// Else case.
				if len(node.child) > 1 {
					err = t.writeNode(w, &node.child[1], ctx)
				}
			}
			if err != nil {
				return err
			}
		}
	case typeCond:
		// Condition node evaluates condition expressions.
		var r bool
		switch {
		case len(node.condHlp) > 0 && node.condLC == lcNone:
			// Condition helper caught (no LC case).
			fn := GetCondFn(byteconv.B2S(node.condHlp))
			if fn == nil {
				err = ErrCondHlpNotFound
				return
			}
			// Prepare arguments list.
			ctx.bufA = ctx.bufA[:0]
			if n := len(node.condHlpArg); n > 0 {
				_ = node.condHlpArg[n-1]
				for i := 0; i < len(node.condHlpArg); i++ {
					arg_ := node.condHlpArg[i]
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
			r, err = t.nodeCmp(node, ctx)
		}
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
		// Evaluate condition.
		if r {
			// True case.
			if len(node.child) > 0 {
				err = t.writeNode(w, &node.child[0], ctx)
			}
		} else {
			// Else case.
			if len(node.child) > 1 {
				err = t.writeNode(w, &node.child[1], ctx)
			}
		}
	case typeCondTrue, typeCondFalse, typeCase, typeDefault:
		// Just walk over child nodes.
		for i := 0; i < len(node.child); i++ {
			ch := &node.child[i]
			err = t.writeNode(w, ch, ctx)
			if err != nil {
				return
			}
		}
	case typeLoopCount:
		// Evaluate counter loops.
		// See Ctx.cloop().
		ctx.brkD = 0
		ctx.cloop(node, t, w)
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
	case typeLoopRange:
		// Evaluate range loops.
		// See Ctx.rloop().
		ctx.brkD = 0
		ctx.rloop(node.loopSrc, node, t, w)
		if ctx.Err != nil {
			err = ctx.Err
			return
		}
	case typeBreak:
		// Break the loop.
		ctx.brkD = node.loopBrkD
		err = ErrBreakLoop
	case typeLBreak:
		// Lazy break the loop.
		ctx.brkD = node.loopBrkD
		err = ErrLBreakLoop
	case typeContinue:
		// Go to next iteration of loop.
		err = ErrContLoop
	case typeSwitch:
		// Switch magic...
		r := false
		if len(node.switchArg) > 0 {
			// Classic switch case.
			for i := 0; i < len(node.child); i++ {
				ch := &node.child[i]
				if ch.typ == typeCase {
					if ch.caseStaticL {
						r = ctx.cmp(node.switchArg, opEq, ch.caseL)
					} else {
						ctx.get(ch.caseL)
						if ctx.Err == nil {
							if err = ctx.BufAcc.StakeOut().WriteX(ctx.bufX).Error(); err != nil {
								return
							}
							r = ctx.cmp(node.switchArg, opEq, ctx.BufAcc.StakedBytes())
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
			for i := 0; i < len(node.child); i++ {
				ch := &node.child[i]
				if ch.typ == typeCase {
					if len(ch.caseHlp) > 0 {
						// Case condition helper caught.
						fn := GetCondFn(byteconv.B2S(ch.caseHlp))
						if fn == nil {
							err = ErrCondHlpNotFound
							return
						}
						// Prepare arguments list.
						ctx.bufA = ctx.bufA[:0]
						if n := len(ch.caseHlpArg); n > 0 {
							_ = ch.caseHlpArg[n-1]
							for j := 0; j < n; j++ {
								arg_ := ch.caseHlpArg[j]
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
			for i := 0; i < len(node.child); i++ {
				ch := &node.child[i]
				if ch.typ == typeDefault {
					err = t.writeNode(w, ch, ctx)
					break
				}
			}
		}
	case typeInclude:
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
	case typeExit:
		// Interrupt template evaluation.
		err = ErrInterrupt
	case typeJsonQ:
		ctx.chJQ = true
	case typeEndJsonQ:
		ctx.chJQ = false
	case typeHtmlE:
		ctx.chHE = true
	case typeEndHtmlE:
		ctx.chHE = false
	case typeUrlEnc:
		ctx.chUE = true
	case typeEndUrlEnc:
		ctx.chUE = false
	default:
		// Unknown node type caught.
		err = ErrUnknownCtl
	}
	return
}

// Evaluate condition expressions.
func (t *Tpl) nodeCmp(node *node, ctx *Ctx) (r bool, err error) {
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
