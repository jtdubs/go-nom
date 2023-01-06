package nom

import (
	"context"
	"testing"
)

type testTracer struct {
	enterCounts, exitCounts map[string]int
}

func newTracer() *testTracer {
	return &testTracer{
		enterCounts: make(map[string]int),
		exitCounts:  make(map[string]int),
	}
}

func (t *testTracer) Enter(_ context.Context, name string, cursor Cursor[rune]) {
	t.enterCounts[name]++
}

func (t *testTracer) Exit(_ context.Context, name string, start, end Cursor[rune], result any, err error) {
	t.exitCounts[name]++
}

func testParseWord(ctx context.Context, start Cursor[rune]) (Cursor[rune], []rune, error) {
	return Trace(Many1(testParseChar))(ctx, start)
}

func testParseChar(ctx context.Context, start Cursor[rune]) (Cursor[rune], rune, error) {
	return Trace(Satisfy(func(r rune) bool { return r >= '0' && r <= '9' }))(ctx, start)
}

func TestTracing(t *testing.T) {
	TraceSupported()

	tracer := newTracer()
	c := NewCursor([]rune("123456"))
	ctx := WithTracing(WithTracer[rune](context.Background(), tracer))

	var wantWord, wantChar int
	for i := 0; i < 10; i = i + 1 {
		if i == 5 {
			ctx = WithoutTracing(ctx)
		}
		testParseWord(ctx, c)
		if i < 5 {
			wantWord = wantWord + 1
			wantChar = wantChar + 7
		}
		if got := tracer.enterCounts["go-nom.testParseWord"]; got != wantWord {
			t.Errorf("enterCounts(nom.testParseWord) = %v, want %v", got, wantWord)
		}
		if got := tracer.enterCounts["go-nom.testParseChar"]; got != wantChar {
			t.Errorf("enterCounts(nom.testParseChar) = %v, want %v", got, wantChar)
		}
		if got := tracer.exitCounts["go-nom.testParseWord"]; got != wantWord {
			t.Errorf("exitCounts(nom.testParseWord) = %v, want %v", got, wantWord)
		}
		if got := tracer.exitCounts["go-nom.testParseChar"]; got != wantChar {
			t.Errorf("exitCounts(nom.testParseChar) = %v, want %v", got, wantChar)
		}
	}
}
