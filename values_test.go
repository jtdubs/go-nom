package nom

import "testing"

func TestAny(t *testing.T) {
	p := Any[rune]
	validate(t, "Any(%q)", p, "\nfoo", 1, '\n', false)
	validate(t, "Any(%q)", p, "Hello", 1, 'H', false)
	validate(t, "Any(%q)", p, "#", 1, '#', false)
	validate(t, "Any(%q)", p, "", 0, rune(0), true)
}

func TestRest(t *testing.T) {
	p := Rest[rune]
	validate(t, "Rest(%q)", p, "Hello", 5, []rune("Hello"), false)
	validate(t, "Rest(%q)", p, "H", 1, []rune("H"), false)
	validate(t, "Rest(%q)", p, "", 0, []rune(""), false)
}

func TestRestLen(t *testing.T) {
	p := RestLen[rune]
	validate(t, "RestLen(%q)", p, "Hello", 5, 5, false)
	validate(t, "RestLen(%q)", p, "H", 1, 1, false)
	validate(t, "RestLen(%q)", p, "", 0, 0, false)
}

func TestEOF(t *testing.T) {
	p := EOF[rune]
	validate(t, "EOF(%q)", p, "Hello", 0, struct{}{}, true)
	validate(t, "EOF(%q)", p, "", 0, struct{}{}, false)
}
