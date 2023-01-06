package runes

import (
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
	return nom.Trace(func(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
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
	return nom.Trace(func(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
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

func EOL(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Preceded(nom.Opt(Rune('\r')), Rune('\n')))(start)
}

func Newline(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return Rune('\n')(start)
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

func Alpha(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsAlpha))(start)
}

func Alpha0(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsAlpha))))(start)
}

func Alpha1(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsAlpha))))(start)
}

func Digit(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsDigit))(start)
}

func Digit0(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsDigit))))(start)
}

func Digit1(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsDigit))))(start)
}

func HexDigit(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsHexDigit))(start)
}

func HexDigit0(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsHexDigit))))(start)
}

func HexDigit1(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsHexDigit))))(start)
}

func OctalDigit(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsOctalDigit))(start)
}

func OctalDigit0(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsOctalDigit))))(start)
}

func OctalDigit1(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsOctalDigit))))(start)
}

func Alphanumeric(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsAlphanumeric))(start)
}

func Alphanumeric0(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsAlphanumeric))))(start)
}

func Alphanumeric1(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsAlphanumeric))))(start)
}

func Space(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsSpace))(start)
}

func Space0(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsSpace))))(start)
}

func Space1(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsSpace))))(start)
}

func Multispace(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsMultispace))(start)
}

func Multispace0(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many0(nom.Satisfy(IsMultispace))))(start)
}

func Multispace1(start nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
	return nom.Trace(Join(nom.Many1(nom.Satisfy(IsMultispace))))(start)
}

func Sign(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return nom.Trace(nom.Satisfy(IsSign))(start)
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
