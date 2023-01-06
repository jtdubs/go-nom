package trace

import (
	"context"
	"testing"

	"github.com/jtdubs/go-nom"
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

func (t *testTracer) Enter(_ context.Context, name string, cursor nom.Cursor[rune]) {
	t.enterCounts[name]++
}

func (t *testTracer) Exit(_ context.Context, name string, start, end nom.Cursor[rune], result any, err error) {
	t.exitCounts[name]++
}

func testParseWord(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], []rune, error) {
	return Trace(nom.Many1(testParseChar))(ctx, start)
}

func testParseChar(ctx context.Context, start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return Trace(nom.Satisfy(func(r rune) bool { return r >= '0' && r <= '9' }))(ctx, start)
}

func TestTracing(t *testing.T) {
	TraceSupported()

	tracer := newTracer()
	c := nom.NewCursor([]rune("123456"))
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
		if got := tracer.enterCounts["trace.testParseWord"]; got != wantWord {
			t.Errorf("enterCounts(trace.testParseWord) = %v, want %v", got, wantWord)
		}
		if got := tracer.enterCounts["trace.testParseChar"]; got != wantChar {
			t.Errorf("enterCounts(trace.testParseChar) = %v, want %v", got, wantChar)
		}
		if got := tracer.exitCounts["trace.testParseWord"]; got != wantWord {
			t.Errorf("exitCounts(trace.testParseWord) = %v, want %v", got, wantWord)
		}
		if got := tracer.exitCounts["trace.testParseChar"]; got != wantChar {
			t.Errorf("exitCounts(trace.testParseChar) = %v, want %v", got, wantChar)
		}
	}
}
