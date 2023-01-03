package cache

import (
	"fmt"
	"testing"

	"github.com/jtdubs/go-nom"
)

func TestCache(t *testing.T) {
	var count int

	f := func(r rune) string { return fmt.Sprintf("Result: %q", r) }

	parseFn := Cache(func(c nom.Cursor[rune]) (nom.Cursor[rune], string, error) {
		count = count + 1
		if c.EOF() {
			return c, "EOF", nil
		}
		result := f(c.Read())
		return c.Next(), result, nil
	})

	for _, msg := range []string{"Hello", "World"} {
		count = 0
		c := nom.NewCursor([]rune(msg))
		for !c.EOF() {
			for i := 0; i < 10; i++ {
				_, got, err := parseFn(c)
				if err != nil {
					t.Errorf("parseFn() returned unexpected err: %v", err)
					return
				}
				if want := f(c.Read()); want != got {
					t.Errorf("parseFn() = %q, want %q", got, want)
					return
				}
			}
			c = c.Next()
		}
		if count != len(msg) {
			t.Errorf("parseFn() count = %q, want %q", count, len(msg))
		}
	}
}
