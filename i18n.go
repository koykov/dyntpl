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
		_ = args[len(args)-1]
		for i := 0; i < len(args); i++ {
			if raw, ok := args[i].(*[]byte); ok && len(*raw) > 0 {
				ctx.repl.AddSolidKV(fastconv.B2S(*raw))
			}
		}
	}

	// Compute the key with preceding locale.
	if len(key) == 0 {
		return nil
	}
	ctx.Buf.Reset().
		WriteStr(ctx.loc).
		WriteByte('.').
		WriteStr(key)

	// Get translation from DB.
	if plural {
		t = db.GetPluralWR(ctx.Buf.String(), def, count, &ctx.repl)
	} else {
		t = db.GetWR(ctx.Buf.String(), def, &ctx.repl)
	}
	ctx.Buf1.Reset().WriteStr(t)
	*buf = &ctx.Buf1

	return nil
}
