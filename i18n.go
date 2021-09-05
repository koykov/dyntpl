package dyntpl

import (
	"bytes"

	"github.com/koykov/fastconv"
	"github.com/koykov/i18n"
)

var (
	defEmpty = []byte(`""`)
)

// Translate label.
func modTranslate(ctx *Ctx, buf *interface{}, _ interface{}, args []interface{}) error {
	return trans(ctx, buf, args, false)
}

// Translate label with plural formula.
func modTranslatePlural(ctx *Ctx, buf *interface{}, _ interface{}, args []interface{}) error {
	return trans(ctx, buf, args, true)
}

// Generic translate function.
func trans(ctx *Ctx, buf *interface{}, args []interface{}, plural bool) error {
	// Check db available.
	if ctx.i18n == nil {
		return nil
	}
	db := (*i18n.DB)(ctx.i18n)

	if len(args) == 0 {
		return ErrModNoArgs
	}

	var (
		key, def, t string
		count       = 1
	)
	// Try to get the key.
	if raw, ok := args[0].(*[]byte); ok && len(*raw) > 0 {
		key = fastconv.B2S(*raw)
	}
	args = args[1:]
	// Try to get the default value.
	if len(args) > 0 {
		if raw, ok := args[0].(*[]byte); ok && len(*raw) > 0 && !bytes.Equal(*raw, defEmpty) {
			def = fastconv.B2S(*raw)
		}
		args = args[1:]
	}
	// Try to get count to use in plural formula.
	if plural {
		if len(args) > 0 {
			if raw, ok := args[0].(int); ok {
				count = raw
			}
			args = args[1:]
		}
	}

	// Collect placeholder replacements.
	ctx.repl.Reset()
	if len(args) > 0 {
		_ = args[len(args)-1]
		for i := 0; i < len(args); i++ {
			if kv, ok := args[i].(*ctxKV); ok {
				ctx.AccBuf.StakeOut().WriteX(kv.v)
				ctx.repl.AddKV(fastconv.B2S(kv.k), ctx.AccBuf.StakedString())
			}
		}
	}

	// Compute the key with preceding locale.
	if len(key) == 0 {
		return nil
	}
	lkey := ctx.AccBuf.StakeOut().WriteStr(ctx.loc).WriteByte('.').WriteStr(key).StakedString()

	// Get translation from DB.
	if plural {
		t = db.GetPluralWR(lkey, def, count, &ctx.repl)
	} else {
		t = db.GetWR(lkey, def, &ctx.repl)
	}
	*buf = ctx.OutBuf.Reset().WriteStr(t)

	return nil
}
