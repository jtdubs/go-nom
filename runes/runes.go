package runes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jtdubs/go-nom"
	"github.com/jtdubs/go-nom/fn"
	"github.com/jtdubs/go-nom/trace"
)

func Rune(want rune) nom.ParseFn[rune, rune] {
	return trace.Trace(func(_ context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
		if start.EOF() {
			return start, rune(0), fmt.Errorf("got %q, want EOF", string(start.Read()))
		}
		if got := start.Read(); got != want {
			return start, rune(0), fmt.Errorf("got %q, want %q", string(got), string(want))
		}
		return start.Next(), want, nil
	})
}

func RuneNoCase(want rune) nom.ParseFn[rune, rune] {
	return trace.Trace(fn.Satisfy(func(got rune) bool {
		return strings.EqualFold(string(want), string(got))
	}))
}

func Tag(tag string) nom.ParseFn[rune, string] {
	runes := make([]nom.ParseFn[rune, rune], len(tag))
	for i, r := range tag {
		runes[i] = Rune(r)
	}
	return trace.Trace(Join(fn.Seq(runes...)))
}

func TagNoCase(tag string) nom.ParseFn[rune, string] {
	runes := make([]nom.ParseFn[rune, rune], len(tag))
	for i, r := range tag {
		runes[i] = RuneNoCase(r)
	}
	return trace.Trace(Join(fn.Seq(runes...)))
}

func OneOf(allowlist string) nom.ParseFn[rune, rune] {
	return trace.Trace(func(_ context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
		if start.EOF() {
			return start, rune(0), errors.New("EOF")
		}
		got := start.Read()
		if !strings.ContainsRune(allowlist, got) {
			return start, rune(0), fmt.Errorf("%q not allowed", got)
		}
		return start.Next(), got, nil
	})
}

func NoneOf(blocklist string) nom.ParseFn[rune, rune] {
	return trace.Trace(func(_ context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
		if start.EOF() {
			return start, rune(0), errors.New("EOF")
		}
		got := start.Read()
		if strings.ContainsRune(blocklist, got) {
			return start, rune(0), fmt.Errorf("%q not allowed", got)
		}
		return start.Next(), got, nil
	})
}

func EOL(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Preceded(fn.Opt(Rune('\r')), Rune('\n')))(ctx, start)
}

func Newline(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(Rune('\n'))(ctx, start)
}

func IsAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func IsAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func IsDigit(r rune) bool {
	return (r >= '0' && r <= '9')
}

func IsHexDigit(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func IsOctalDigit(r rune) bool {
	return (r >= '0' && r <= '7')
}

func IsSpace(r rune) bool {
	return r == ' ' || r == '\t'
}

func IsMultispace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func IsSign(r rune) bool {
	return r == '+' || r == '-'
}

func Alpha(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsAlpha))(ctx, start)
}

func Alpha0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many0(fn.Satisfy(IsAlpha))))(ctx, start)
}

func Alpha1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many1(fn.Satisfy(IsAlpha))))(ctx, start)
}

func Digit(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsDigit))(ctx, start)
}

func Digit0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many0(fn.Satisfy(IsDigit))))(ctx, start)
}

func Digit1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many1(fn.Satisfy(IsDigit))))(ctx, start)
}

func HexDigit(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsHexDigit))(ctx, start)
}

func HexDigit0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many0(fn.Satisfy(IsHexDigit))))(ctx, start)
}

func HexDigit1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many1(fn.Satisfy(IsHexDigit))))(ctx, start)
}

func OctalDigit(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsOctalDigit))(ctx, start)
}

func OctalDigit0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many0(fn.Satisfy(IsOctalDigit))))(ctx, start)
}

func OctalDigit1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many1(fn.Satisfy(IsOctalDigit))))(ctx, start)
}

func Alphanumeric(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsAlphanumeric))(ctx, start)
}

func Alphanumeric0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many0(fn.Satisfy(IsAlphanumeric))))(ctx, start)
}

func Alphanumeric1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many1(fn.Satisfy(IsAlphanumeric))))(ctx, start)
}

func Space(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsSpace))(ctx, start)
}

func Space0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many0(fn.Satisfy(IsSpace))))(ctx, start)
}

func Space1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many1(fn.Satisfy(IsSpace))))(ctx, start)
}

func Multispace(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsMultispace))(ctx, start)
}

func Multispace0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many0(fn.Satisfy(IsMultispace))))(ctx, start)
}

func Multispace1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return trace.Trace(Join(fn.Many1(fn.Satisfy(IsMultispace))))(ctx, start)
}

func Sign(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return trace.Trace(fn.Satisfy(IsSign))(ctx, start)
}

func Phrase[T any](ps ...nom.ParseFn[rune, T]) nom.ParseFn[rune, []T] {
	var parts []nom.ParseFn[rune, T]
	for _, p := range ps {
		parts = append(parts, fn.Preceded(Space0, p))
	}
	return trace.Trace(fn.Seq(parts...))
}

func Surrounded[T, U, V any](left nom.ParseFn[rune, T], right nom.ParseFn[rune, U], middle nom.ParseFn[rune, V]) nom.ParseFn[rune, V] {
	return trace.Trace(fn.Surrounded(
		fn.Preceded(Space0, left),
		fn.Preceded(Space0, right),
		fn.Preceded(Space0, middle),
	))
}

func SurroundedBy[T any](left, right rune, middle nom.ParseFn[rune, T]) nom.ParseFn[rune, T] {
	return trace.Trace(Surrounded(Rune(left), Rune(right), middle))
}

func Recognize[T any](p nom.ParseFn[rune, T]) nom.ParseFn[rune, string] {
	return trace.Trace(Join(fn.Recognize(p)))
}

func Concat(p nom.ParseFn[rune, []string]) nom.ParseFn[rune, string] {
	return trace.Trace(fn.Map(p, func(ss []string) string {
		var result string
		for _, s := range ss {
			result = result + s
		}
		return result
	}))
}

func Join(p nom.ParseFn[rune, []rune]) nom.ParseFn[rune, string] {
	return trace.Trace(fn.Map(p, func(rs []rune) string { return string(rs) }))
}

func Cons(p nom.ParseFn[rune, rune], ps nom.ParseFn[rune, string]) nom.ParseFn[rune, string] {
	return trace.Trace(Concat(fn.Seq(fn.Map(p, func(r rune) string { return string(r) }), ps)))
}

func Cursor(s string) nom.Cursor[rune] {
	return nom.NewCursor([]rune(s))
}
