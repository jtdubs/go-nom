package runes

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jtdubs/go-nom"
)

func validate[T any](t *testing.T, name string, p nom.ParseFn[rune, T], in string, wantPosition int, wantResult T, wantError bool) {
	t.Helper()

	name = fmt.Sprintf(name, in)
	inCursor := Cursor(in)
	gotCursor, gotResult, err := p(context.Background(), inCursor)
	if gotCursor.Position() != wantPosition {
		t.Errorf("%v(%v) cursor = %v, want %v", name, inCursor, gotCursor.Position(), wantPosition)
		return
	}
	if diff := cmp.Diff(wantResult, gotResult, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("%v(%v) result unexpected diff (-want +got):\n%v\n", name, inCursor, diff)
		return
	}
	if gotError := (err != nil); gotError != wantError {
		if wantError {
			t.Errorf("%v(%v) = '%v', want error", name, inCursor, gotResult)
		} else {
			t.Errorf("%v(%v) unexpected error: %v", name, inCursor, err)
		}
		return
	}
}

func TestRune(t *testing.T) {
	p := Rune('H')
	validate(t, "Rune(%q)", p, "Hello", 1, 'H', false)
	validate(t, "Rune(%q)", p, "hello", 0, rune(0), true)
	validate(t, "Rune(%q)", p, "", 0, rune(0), true)
}

func TestRuneNoCase(t *testing.T) {
	p := RuneNoCase('H')
	validate(t, "RuneNoCase(%q)", p, "Hello", 1, 'H', false)
	validate(t, "RuneNoCase(%q)", p, "hello", 1, 'h', false)
	validate(t, "RuneNoCase(%q)", p, "Jello", 0, rune(0), true)
	validate(t, "RuneNoCase(%q)", p, "", 0, rune(0), true)
}

func TestTag(t *testing.T) {
	p := Tag("Hello")
	validate(t, "Tag(%q)", p, "Hello World", 5, "Hello", false)
	validate(t, "Tag(%q)", p, "Hello ", 5, "Hello", false)
	validate(t, "Tag(%q)", p, "Hello", 5, "Hello", false)
	validate(t, "Tag(%q)", p, "Hellf", 0, "", true)
	validate(t, "Tag(%q)", p, "Hell", 0, "", true)
	validate(t, "Tag(%q)", p, "H", 0, "", true)
	validate(t, "Tag(%q)", p, "", 0, "", true)
}

func TestTagNoCase(t *testing.T) {
	p := TagNoCase("Hello")
	validate(t, "TagNoCase(%q)", p, "Hello World", 5, "Hello", false)
	validate(t, "TagNoCase(%q)", p, "Hello ", 5, "Hello", false)
	validate(t, "TagNoCase(%q)", p, "Hello", 5, "Hello", false)
	validate(t, "TagNoCase(%q)", p, "Hellf", 0, "", true)
	validate(t, "TagNoCase(%q)", p, "Hell", 0, "", true)
	validate(t, "TagNoCase(%q)", p, "H", 0, "", true)
	validate(t, "TagNoCase(%q)", p, "", 0, "", true)
	validate(t, "TagNoCase(%q)", p, "HELLO", 5, "HELLO", false)
	validate(t, "TagNoCase(%q)", p, "hello", 5, "hello", false)
}

func TestOneOf(t *testing.T) {
	p := OneOf("HJ")
	validate(t, "OneOf(%q)", p, "Hello", 1, 'H', false)
	validate(t, "OneOf(%q)", p, "Jello", 1, 'J', false)
	validate(t, "OneOf(%q)", p, "Mello", 0, rune(0), true)
	validate(t, "OneOf(%q)", p, "", 0, rune(0), true)
}

func TestNoneOf(t *testing.T) {
	p := NoneOf("MK")
	validate(t, "NoneOf(%q)", p, "Hello", 1, 'H', false)
	validate(t, "NoneOf(%q)", p, "Jello", 1, 'J', false)
	validate(t, "NoneOf(%q)", p, "Mello", 0, rune(0), true)
	validate(t, "NoneOf(%q)", p, "Kello", 0, rune(0), true)
	validate(t, "NoneOf(%q)", p, "", 0, rune(0), true)
}

func TestEOL(t *testing.T) {
	p := EOL
	validate(t, "EOL(%q)", p, "\r\n", 2, '\n', false)
	validate(t, "EOL(%q)", p, "\n", 1, '\n', false)
	validate(t, "EOL(%q)", p, "\r", 0, rune(0), true)
	validate(t, "EOL(%q)", p, "abc", 0, rune(0), true)
	validate(t, "EOL(%q)", p, "", 0, rune(0), true)
}

func TestNewline(t *testing.T) {
	p := Newline
	validate(t, "Newline(%q)", p, "\n", 1, '\n', false)
	validate(t, "Newline(%q)", p, "\r\n", 0, rune(0), true)
	validate(t, "Newline(%q)", p, "\r", 0, rune(0), true)
	validate(t, "Newline(%q)", p, "abc", 0, rune(0), true)
	validate(t, "Newline(%q)", p, "", 0, rune(0), true)
}

func TestAlpha(t *testing.T) {
	p := Alpha
	validate(t, "Alpha(%q)", p, "Hello World", 1, 'H', false)
	validate(t, "Alpha(%q)", p, "H123", 1, 'H', false)
	validate(t, "Alpha(%q)", p, "123Hello", 0, rune(0), true)
	validate(t, "Alpha(%q)", p, "", 0, rune(0), true)
}

func TestAlpha0(t *testing.T) {
	p := Alpha0
	validate(t, "Alpha0(%q)", p, "Hello World", 5, "Hello", false)
	validate(t, "Alpha0(%q)", p, "Hello1World", 5, "Hello", false)
	validate(t, "Alpha0(%q)", p, "H123", 1, "H", false)
	validate(t, "Alpha0(%q)", p, "123Hello", 0, "", false)
	validate(t, "Alpha0(%q)", p, "", 0, "", false)
}

func TestAlpha1(t *testing.T) {
	p := Alpha1
	validate(t, "Alpha1(%q)", p, "Hello World", 5, "Hello", false)
	validate(t, "Alpha1(%q)", p, "Hello1World", 5, "Hello", false)
	validate(t, "Alpha1(%q)", p, "H123", 1, "H", false)
	validate(t, "Alpha1(%q)", p, "123Hello", 0, "", true)
	validate(t, "Alpha1(%q)", p, "", 0, "", true)
}

func TestDigit(t *testing.T) {
	p := Digit
	validate(t, "Digit(%q)", p, "123 hello", 1, '1', false)
	validate(t, "Digit(%q)", p, "1g", 1, '1', false)
	validate(t, "Digit(%q)", p, "hello", 0, rune(0), true)
	validate(t, "Digit(%q)", p, "", 0, rune(0), true)
}

func TestDigit0(t *testing.T) {
	p := Digit0
	validate(t, "Digit0(%q)", p, "123 hello", 3, "123", false)
	validate(t, "Digit0(%q)", p, "1g", 1, "1", false)
	validate(t, "Digit0(%q)", p, "hello", 0, "", false)
	validate(t, "Digit0(%q)", p, "", 0, "", false)
}

func TestDigit1(t *testing.T) {
	p := Digit1
	validate(t, "Digit1(%q)", p, "123 hello", 3, "123", false)
	validate(t, "Digit1(%q)", p, "1g", 1, "1", false)
	validate(t, "Digit1(%q)", p, "hello", 0, "", true)
	validate(t, "Digit1(%q)", p, "", 0, "", true)
}

func TestHexDigit(t *testing.T) {
	p := HexDigit
	validate(t, "HexDigit(%q)", p, "123 hello", 1, '1', false)
	validate(t, "HexDigit(%q)", p, "1g", 1, '1', false)
	validate(t, "HexDigit(%q)", p, "f1c23q", 1, 'f', false)
	validate(t, "HexDigit(%q)", p, "jello", 0, rune(0), true)
	validate(t, "HexDigit(%q)", p, "", 0, rune(0), true)
}

func TestHexDigit0(t *testing.T) {
	p := HexDigit0
	validate(t, "HexDigit0(%q)", p, "123 hello", 3, "123", false)
	validate(t, "HexDigit0(%q)", p, "1g", 1, "1", false)
	validate(t, "HexDigit0(%q)", p, "1fc23q", 5, "1fc23", false)
	validate(t, "HexDigit0(%q)", p, "jello", 0, "", false)
	validate(t, "HexDigit0(%q)", p, "", 0, "", false)
}

func TestHexDigit1(t *testing.T) {
	p := HexDigit1
	validate(t, "HexDigit1(%q)", p, "123 hello", 3, "123", false)
	validate(t, "HexDigit1(%q)", p, "1g", 1, "1", false)
	validate(t, "HexDigit1(%q)", p, "1fc23q", 5, "1fc23", false)
	validate(t, "HexDigit1(%q)", p, "jello", 0, "", true)
	validate(t, "HexDigit1(%q)", p, "", 0, "", true)
}

func TestOctalDigit(t *testing.T) {
	p := OctalDigit
	validate(t, "OctalDigit(%q)", p, "123 hello", 1, '1', false)
	validate(t, "OctalDigit(%q)", p, "18", 1, '1', false)
	validate(t, "OctalDigit(%q)", p, "jello", 0, rune(0), true)
	validate(t, "OctalDigit(%q)", p, "", 0, rune(0), true)
}

func TestOctalDigit0(t *testing.T) {
	p := OctalDigit0
	validate(t, "OctalDigit0(%q)", p, "123 hello", 3, "123", false)
	validate(t, "OctalDigit0(%q)", p, "18", 1, "1", false)
	validate(t, "OctalDigit0(%q)", p, "jello", 0, "", false)
	validate(t, "OctalDigit0(%q)", p, "", 0, "", false)
}

func TestOctalDigit1(t *testing.T) {
	p := OctalDigit1
	validate(t, "OctalDigit1(%q)", p, "123 hello", 3, "123", false)
	validate(t, "OctalDigit1(%q)", p, "18", 1, "1", false)
	validate(t, "OctalDigit1(%q)", p, "jello", 0, "", true)
	validate(t, "OctalDigit1(%q)", p, "", 0, "", true)
}

func TestAlphanumeric0(t *testing.T) {
	p := Alphanumeric0
	validate(t, "Alphanumeric0(%q)", p, "hello123world", 13, "hello123world", false)
	validate(t, "Alphanumeric0(%q)", p, "h", 1, "h", false)
	validate(t, "Alphanumeric0(%q)", p, "", 0, "", false)
}

func TestAlphanumeric(t *testing.T) {
	p := Alphanumeric
	validate(t, "Alphanumeric(%q)", p, "hello123world", 1, 'h', false)
	validate(t, "Alphanumeric(%q)", p, "h", 1, 'h', false)
	validate(t, "Alphanumeric(%q)", p, "", 0, rune(0), true)
}

func TestAlphanumeric1(t *testing.T) {
	p := Alphanumeric1
	validate(t, "Alphanumeric1(%q)", p, "hello123world", 13, "hello123world", false)
	validate(t, "Alphanumeric1(%q)", p, "h", 1, "h", false)
	validate(t, "Alphanumeric1(%q)", p, "", 0, "", true)
}

func TestSpace(t *testing.T) {
	p := Space
	validate(t, "Space(%q)", p, " \t  hi", 1, ' ', false)
	validate(t, "Space(%q)", p, "\t hi", 1, '\t', false)
	validate(t, "Space(%q)", p, "hi", 0, rune(0), true)
	validate(t, "Space(%q)", p, "", 0, rune(0), true)
}

func TestSpace0(t *testing.T) {
	p := Space0
	validate(t, "Space0(%q)", p, " \t  hi", 4, " \t  ", false)
	validate(t, "Space0(%q)", p, " hi", 1, " ", false)
	validate(t, "Space0(%q)", p, "hi", 0, "", false)
	validate(t, "Space0(%q)", p, "", 0, "", false)
}

func TestSpace1(t *testing.T) {
	p := Space1
	validate(t, "Space1(%q)", p, " \t  hi", 4, " \t  ", false)
	validate(t, "Space1(%q)", p, " hi", 1, " ", false)
	validate(t, "Space1(%q)", p, "hi", 0, "", true)
	validate(t, "Space1(%q)", p, "", 0, "", true)
}

func TestMultispace(t *testing.T) {
	p := Multispace
	validate(t, "Multispace(%q)", p, " \t  \r\nhi", 1, ' ', false)
	validate(t, "Multispace(%q)", p, "\r hi", 1, '\r', false)
	validate(t, "Multispace(%q)", p, "hi", 0, rune(0), true)
	validate(t, "Multispace(%q)", p, "", 0, rune(0), true)
}

func TestMultispace0(t *testing.T) {
	p := Multispace0
	validate(t, "Multispace0(%q)", p, " \t  \r\nhi", 6, " \t  \r\n", false)
	validate(t, "Multispace0(%q)", p, "\r hi", 2, "\r ", false)
	validate(t, "Multispace0(%q)", p, "hi", 0, "", false)
	validate(t, "Multispace0(%q)", p, "", 0, "", false)
}

func TestMultispace1(t *testing.T) {
	p := Multispace1
	validate(t, "Multispace1(%q)", p, " \t  \r\nhi", 6, " \t  \r\n", false)
	validate(t, "Multispace1(%q)", p, "\r hi", 2, "\r ", false)
	validate(t, "Multispace1(%q)", p, "hi", 0, "", true)
	validate(t, "Multispace1(%q)", p, "", 0, "", true)
}

func TestSign(t *testing.T) {
	p := Sign
	validate(t, "Sign(%q)", p, "+foo", 1, '+', false)
	validate(t, "Sign(%q)", p, "-foo", 1, '-', false)
	validate(t, "Sign(%q)", p, "f+oo", 0, rune(0), true)
	validate(t, "Sign(%q)", p, "", 0, rune(0), true)
}

func TestPhrase(t *testing.T) {
	p := Phrase(Alpha1, Alpha1)
	validate(t, "Alpha1(%q)", p, "hello world", 11, []string{"hello", "world"}, false)
	validate(t, "Alpha1(%q)", p, "   hello world", 14, []string{"hello", "world"}, false)
	validate(t, "Alpha1(%q)", p, "hello \t world  ", 13, []string{"hello", "world"}, false)
	validate(t, "Alpha1(%q)", p, "hello \t ", 0, nil, true)
}

func TestSurrounded(t *testing.T) {
	p := SurroundedBy('(', ')', Alpha1)
	validate(t, "SurroundedBy(%q)", p, " ( hi ) ", 7, "hi", false)
	validate(t, "SurroundedBy(%q)", p, " ( hi | ", 0, "", true)
	validate(t, "SurroundedBy(%q)", p, " | hi ) ", 0, "", true)
	validate(t, "SurroundedBy(%q)", p, " ( 123 ) ", 0, "", true)
}

func TestRecognize(t *testing.T) {
	p := Recognize(nom.Discard(Alpha1))
	validate(t, "Recognize(%q)", p, "Hello world", 5, "Hello", false)
	validate(t, "Recognize(%q)", p, "H", 1, "H", false)
	validate(t, "Recognize(%q)", p, "", 0, "", true)
}

func TestConcat(t *testing.T) {
	p := Concat(nom.Seq(Alpha1, Space1, Alpha1))
	validate(t, "Concat(%q)", p, "Hello world", 11, "Hello world", false)
	validate(t, "Concat(%q)", p, "H", 0, "", true)
	validate(t, "Concat(%q)", p, "", 0, "", true)
}

func TestCons(t *testing.T) {
	p := Cons(Alpha, Digit0)
	validate(t, "Concat(%q)", p, "a123", 4, "a123", false)
	validate(t, "Concat(%q)", p, "a", 1, "a", false)
	validate(t, "Concat(%q)", p, "abc", 1, "a", false)
	validate(t, "Concat(%q)", p, "123", 0, "", true)
}
