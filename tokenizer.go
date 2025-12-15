package dyntpl

import (
	"github.com/koykov/byteseq"
	"github.com/koykov/simd/indextoken"
)

// Tokenize splits s to tokens and appends them to dst.
// Example: "disk.partitions[0].size@gb" -> ["disk", "partitions", "0", "size", "gb"]
func Tokenize[T byteseq.Q](dst []T, s T) []T {
	var tkn indextoken.Tokenizer[T]
	for {
		t := tkn.Next(s)
		if len(t) == 0 {
			return dst
		}
		dst = append(dst, t)
	}
}

func tokenize(dst []string, s string) []string {
	var tkn indextoken.Tokenizer[string]
	tkn.KeepSquareBrackets()
	for {
		t := tkn.Next(s)
		if len(t) == 0 {
			return dst
		}
		dst = append(dst, t)
	}
}

var _ = Tokenize[string]
