package nom

import (
	"testing"
)

func TestSeq(t *testing.T) {
	p := Seq(Expect('H'), Expect('e'))
	validate(t, "Seq(%q)", p, "Hello", 2, []rune("He"), false)
	validate(t, "Seq(%q)", p, "Gello", 0, []rune(""), true)
	validate(t, "Seq(%q)", p, "Hfllo", 0, []rune(""), true)
	validate(t, "Seq(%q)", p, "H", 0, []rune(""), true)
	validate(t, "Seq(%q)", p, "", 0, []rune(""), true)
}

func TestSurrounded(t *testing.T) {
	p := Surrounded(Expect('('), Expect(')'), Expect('H'))
	validate(t, "Surrounded(%q)", p, "(H)", 3, 'H', false)
	validate(t, "Surrounded(%q)", p, "(I)", 0, rune(0), true)
	validate(t, "Surrounded(%q)", p, "H", 0, rune(0), true)
	validate(t, "Surrounded(%q)", p, "H)", 0, rune(0), true)
	validate(t, "Surrounded(%q)", p, "(H", 0, rune(0), true)
}

func TestPreceded(t *testing.T) {
	p := Preceded(Expect('H'), Expect('e'))
	validate(t, "Preceded(%q)", p, "Hello", 2, 'e', false)
	validate(t, "Preceded(%q)", p, "He", 2, 'e', false)
	validate(t, "Preceded(%q)", p, "Jello", 0, rune(0), true)
	validate(t, "Preceded(%q)", p, "Hfllo", 0, rune(0), true)
	validate(t, "Preceded(%q)", p, "Hf", 0, rune(0), true)
	validate(t, "Preceded(%q)", p, "H", 0, rune(0), true)
}

func TestTerminated(t *testing.T) {
	p := Terminated(Expect('H'), Expect('e'))
	validate(t, "Terminated(%q)", p, "Hello", 2, 'H', false)
	validate(t, "Terminated(%q)", p, "He", 2, 'H', false)
	validate(t, "Terminated(%q)", p, "Jello", 0, rune(0), true)
	validate(t, "Terminated(%q)", p, "Hfllo", 0, rune(0), true)
	validate(t, "Terminated(%q)", p, "Hf", 0, rune(0), true)
	validate(t, "Terminated(%q)", p, "H", 0, rune(0), true)
}
