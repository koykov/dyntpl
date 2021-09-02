package dyntpl

import (
	"bytes"

	"github.com/koykov/fastconv"
	"github.com/koykov/i18n"
	"github.com/koykov/x2bytes"
)

var (
	defEmpty = []byte(`""`)
)

// Translate label.
func modTranslate(ctx *Ctx, buf *interface{}, _ interface{}, args []interface{}) (err error) {
	return trans(ctx, buf, args, false)
}

// Translate label with plural formula.
func modTranslatePlural(ctx *Ctx, buf *interface{}, _ interface{}, args []interface{}) (err error) {
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
		var err error
		bufArgs := ctx.GetByteBuf()
		bufArgs.Reset()
		_ = args[len(args)-1]
		for i := 0; i < len(args); i++ {
			if kv, ok := args[i].(*ctxKV); ok {
				off := bufArgs.Len()
				if *bufArgs, err = x2bytes.ToBytes(*bufArgs, kv.v); err == nil {
					ctx.repl.AddKV(fastconv.B2S(kv.k), fastconv.B2S((*bufArgs)[off:bufArgs.Len()]))
				}
			}
		}
	}

	// Compute the key with preceding locale.
	if len(key) == 0 {
		return nil
	}
	bufKey := ctx.GetByteBuf().Reset().WriteStr(ctx.loc).WriteByte('.').WriteStr(key)

	// Get translation from DB.
	if plural {
		t = db.GetPluralWR(bufKey.String(), def, count, &ctx.repl)
	} else {
		t = db.GetWR(bufKey.String(), def, &ctx.repl)
	}
	bufResult := ctx.GetByteBuf().Reset().WriteStr(t)
	*buf = bufResult

	return nil
}
