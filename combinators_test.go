package nom

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func validate[T any](t *testing.T, name string, p ParseFn[rune, T], in string, wantPosition int, wantResult T, wantError bool) bool {
	t.Helper()

	name = fmt.Sprintf(name, in)
	inCursor := NewCursor([]rune(in))
	gotCursor, gotResult, err := p(context.Background(), inCursor)
	if gotCursor.Position() != wantPosition {
		t.Errorf("%v(%q) cursor = %v, want %v", name, inCursor, gotCursor.Position(), wantPosition)
		return false
	}
	if diff := cmp.Diff(wantResult, gotResult, cmpopts.EquateEmpty()); diff != "" {
		t.Errorf("%v(%v) result unexpected diff (-want +got):\n%v\n", name, inCursor, diff)
		return false
	}
	if gotError := (err != nil); gotError != wantError {
		if wantError {
			t.Errorf("%v(%q) = '%v', want error", name, inCursor, gotResult)
		} else {
			t.Errorf("%v(%q) unexpected error: %v", name, inCursor, err)
		}
		return false
	}
	return true
}

func validateBind[T any](t *testing.T, name string, p ParseFn[rune, T], in string, wantPosition int, wantResult T, wantError bool, got *rune, want rune) bool {
	t.Helper()
	*got = rune(0)
	if !validate(t, name, p, in, wantPosition, wantResult, wantError) {
		return false
	}
	if *got != want {
		name = fmt.Sprintf(name, in)
		t.Errorf("%v bound %q, want %q", name, *got, want)
		return false
	}
	return true
}

func TestAlt(t *testing.T) {
	p := Alt(Expect('H'), Expect('I'))
	validate(t, "Alt(%q)", p, "Hello", 1, 'H', false)
	validate(t, "Alt(%q)", p, "Iello", 1, 'I', false)
	validate(t, "Alt(%q)", p, "Jello", 0, rune(0), true)
	validate(t, "Alt(%q)", p, "", 0, rune(0), true)
}

func TestMap(t *testing.T) {
	p := Map(Seq(Expect('H'), Expect('e')), func(rs []rune) int { return len(rs) })
	validate(t, "Map(%q)", p, "Hello", 2, 2, false)
	validate(t, "Map(%q)", p, "Hillo", 0, 0, true)
}

func TestOpt(t *testing.T) {
	p := Opt(Expect('H'))
	validate(t, "Opt(%q)", p, "Hello", 1, 'H', false)
	validate(t, "Opt(%q)", p, "Jello", 0, rune(0), false)
	validate(t, "Opt(%q)", p, "", 0, rune(0), false)
}

func TestCondTrue(t *testing.T) {
	p := Cond(true, Expect('H'))
	validate(t, "Cond(true, %q)", p, "Hello", 1, 'H', false)
	validate(t, "Cond(true, %q)", p, "Jello", 0, rune(0), true)
	validate(t, "Cond(true, %q)", p, "", 0, rune(0), true)
}

func TestCondFalse(t *testing.T) {
	p := Cond(false, Expect('H'))
	validate(t, "Cond(false, %q)", p, "Hello", 0, rune(0), false)
	validate(t, "Cond(false, %q)", p, "Jello", 0, rune(0), false)
	validate(t, "Cond(false, %q)", p, "", 0, rune(0), false)
}

func TestPeek(t *testing.T) {
	p := Peek(Expect('H'))
	validate(t, "Peek(%q)", p, "Hello", 0, 'H', false)
	validate(t, "Peek(%q)", p, "Jello", 0, rune(0), true)
	validate(t, "Peek(%q)", p, "", 0, rune(0), true)
}

func TestVerify(t *testing.T) {
	p := Verify(Alt(Expect('H'), Expect('I')), func(r rune) bool { return r == 'H' })
	validate(t, "Verify(%q)", p, "Hello", 1, 'H', false)
	validate(t, "Verify(%q)", p, "Iello", 0, rune(0), true)
	validate(t, "Verify(%q)", p, "Jello", 0, rune(0), true)
	validate(t, "Verify(%q)", p, "", 0, rune(0), true)
}

func TestValue(t *testing.T) {
	p := Value(42, Expect('H'))
	validate(t, "Value(%q)", p, "Hello", 1, 42, false)
	validate(t, "Value(%q)", p, "Jello", 0, 0, true)
	validate(t, "Value(%q)", p, "", 0, 0, true)
}

func TestNot(t *testing.T) {
	p := Not(Expect('H'))
	validate(t, "Not(%q)", p, "Hello", 0, rune(0), true)
	validate(t, "Not(%q)", p, "Jello", 0, rune(0), false)
	validate(t, "Not(%q)", p, "", 0, rune(0), false)
}

func TestSuccess(t *testing.T) {
	p := Success[rune](42)
	validate(t, "Success(%q)", p, "Hello", 0, 42, false)
	validate(t, "Success(%q)", p, "Jello", 0, 42, false)
	validate(t, "Success(%q)", p, "", 0, 42, false)
}

func TestFailure(t *testing.T) {
	err := errors.New("oops")
	p := Failure[rune, int](err)
	validate(t, "Failure(%q)", p, "Hello", 0, 0, true)
	validate(t, "Failure(%q)", p, "Jello", 0, 0, true)
	validate(t, "Failure(%q)", p, "", 0, 0, true)
}

func TestRecognize(t *testing.T) {
	p := Recognize(Seq(Expects([]rune("Hel")), Expects([]rune("lo"))))
	validate(t, "Recognize(%q)", p, "Hello World", 5, []rune("Hello"), false)
	validate(t, "Recognize(%q)", p, "Hellf World", 0, []rune(""), true)
	validate(t, "Recognize(%q)", p, "Hello ", 5, []rune("Hello"), false)
	validate(t, "Recognize(%q)", p, "Hellf", 0, []rune(""), true)
	validate(t, "Recognize(%q)", p, "Hillo ", 0, []rune(""), true)
	validate(t, "Recognize(%q)", p, "Hell", 0, []rune(""), true)
	validate(t, "Recognize(%q)", p, "H", 0, []rune(""), true)
	validate(t, "Recognize(%q)", p, "J", 0, []rune(""), true)
	validate(t, "Recognize(%q)", p, "", 0, []rune(""), true)
}

func TestBind(t *testing.T) {
	var got rune
	p := Bind(&got, Alt(Expect('H'), Expect('J')))
	validateBind(t, "Bind(%q)", p, "Hello", 1, struct{}{}, false, &got, 'H')
	validateBind(t, "Bind(%q)", p, "Jello", 1, struct{}{}, false, &got, 'J')
	validateBind(t, "Bind(%q)", p, "Cello", 0, struct{}{}, true, &got, rune(0))
	validateBind(t, "Bind(%q)", p, "", 0, struct{}{}, true, &got, rune(0))
}

func TestSatisfy(t *testing.T) {
	p := Satisfy(func(r rune) bool { return r == 'H' })
	validate(t, "Satisfy(%q)", p, "Hello", 1, 'H', false)
	validate(t, "Satisfy(%q)", p, "Jello", 0, rune(0), true)
	validate(t, "Satisfy(%q)", p, "H", 1, 'H', false)
	validate(t, "Satisfy(%q)", p, "J", 0, rune(0), true)
	validate(t, "Satisfy(%q)", p, "", 0, rune(0), true)
}

func TestDiscard(t *testing.T) {
	p := Discard(Expect('H'))
	validate(t, "Discard(%q)", p, "Hello", 1, struct{}{}, false)
	validate(t, "Discard(%q)", p, "Jello", 0, struct{}{}, true)
	validate(t, "Discard(%q)", p, "", 0, struct{}{}, true)
}
