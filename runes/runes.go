package runes

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jtdubs/go-nom"
)

func Rune(want rune) nom.ParseFn[rune, rune] {
	return nom.Trace(nom.Expect(want))
}

func RuneNoCase(want rune) nom.ParseFn[rune, rune] {
	return nom.Trace(nom.Satisfy(func(got rune) bool {
		return strings.EqualFold(string(want), string(got))
	}))
}

func Tag(tag string) nom.ParseFn[rune, string] {
	runes := make([]nom.ParseFn[rune, rune], len(tag))
	for i, r := range tag {
		runes[i] = Rune(r)
	}
	return Join(nom.Seq(runes...))
}

func TagNoCase(tag string) nom.ParseFn[rune, string] {
	runes := make([]nom.ParseFn[rune, rune], len(tag))
	for i, r := range tag {
		runes[i] = RuneNoCase(r)
	}
	return Join(nom.Seq(runes...))
}

func OneOf(allowlist string) nom.ParseFn[rune, rune] {
	return nom.Trace(func(_ context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
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
	return nom.Trace(func(_ context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
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
	return nom.Trace(nom.Preceded(nom.Opt(Rune('\r')), Rune('\n')))(ctx, start)
}

func Newline(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return Rune('\n')(ctx, start)
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
	return nom.Trace(nom.Satisfy(IsAlpha))(ctx, start)
}

func Alpha0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsAlpha))))(ctx, start)
}

func Alpha1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsAlpha))))(ctx, start)
}

func Digit(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsDigit))(ctx, start)
}

func Digit0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsDigit))))(ctx, start)
}

func Digit1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsDigit))))(ctx, start)
}

func HexDigit(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsHexDigit))(ctx, start)
}

func HexDigit0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsHexDigit))))(ctx, start)
}

func HexDigit1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsHexDigit))))(ctx, start)
}

func OctalDigit(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsOctalDigit))(ctx, start)
}

func OctalDigit0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsOctalDigit))))(ctx, start)
}

func OctalDigit1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsOctalDigit))))(ctx, start)
}

func Alphanumeric(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsAlphanumeric))(ctx, start)
}

func Alphanumeric0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsAlphanumeric))))(ctx, start)
}

func Alphanumeric1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsAlphanumeric))))(ctx, start)
}

func Space(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsSpace))(ctx, start)
}

func Space0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsSpace))))(ctx, start)
}

func Space1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsSpace))))(ctx, start)
}

func Multispace(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsMultispace))(ctx, start)
}

func Multispace0(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsMultispace))))(ctx, start)
}

func Multispace1(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsMultispace))))(ctx, start)
}

func Sign(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsSign))(ctx, start)
}

func Phrase[T any](ps ...nom.ParseFn[rune, T]) nom.ParseFn[rune, []T] {
	var parts []nom.ParseFn[rune, T]
	for _, p := range ps {
		parts = append(parts, nom.Preceded(Space0, p))
	}
	return nom.Trace(nom.Seq(parts...))
}

func Surrounded[T, U, V any](left nom.ParseFn[rune, T], right nom.ParseFn[rune, U], middle nom.ParseFn[rune, V]) nom.ParseFn[rune, V] {
	return nom.Trace(nom.Surrounded(
		nom.Preceded(Space0, left),
		nom.Preceded(Space0, right),
		nom.Preceded(Space0, middle),
	))
}

func SurroundedBy[T any](left, right rune, middle nom.ParseFn[rune, T]) nom.ParseFn[rune, T] {
	return nom.Trace(Surrounded(Rune(left), Rune(right), middle))
}

func Recognize[T any](p nom.ParseFn[rune, T]) nom.ParseFn[rune, string] {
	return Join(nom.Recognize(p))
}

func Concat(p nom.ParseFn[rune, []string]) nom.ParseFn[rune, string] {
	return nom.Map(p, func(ss []string) string {
		var result string
		for _, s := range ss {
			result = result + s
		}
		return result
	})
}

func Join(p nom.ParseFn[rune, []rune]) nom.ParseFn[rune, string] {
	return nom.Map(p, func(rs []rune) string { return string(rs) })
}

func Cons(p nom.ParseFn[rune, rune], ps nom.ParseFn[rune, string]) nom.ParseFn[rune, string] {
	return Concat(nom.Seq(nom.Map(p, func(r rune) string { return string(r) }), ps))
}

func Cursor(s string) nom.Cursor[rune] {
	return nom.NewCursor([]rune(s))
}
