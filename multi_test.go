package nom

import (
	"testing"
)

func TestMany0(t *testing.T) {
	p := Many0(Expect('H'))
	validate(t, "Many0(%q)", p, "HHHHH", 5, []rune("HHHHH"), false)
	validate(t, "Many0(%q)", p, "HHHHHJ", 5, []rune("HHHHH"), false)
	validate(t, "Many0(%q)", p, "H", 1, []rune("H"), false)
	validate(t, "Many0(%q)", p, "HJ", 1, []rune("H"), false)
	validate(t, "Many0(%q)", p, "", 0, []rune(""), false)
	validate(t, "Many0(%q)", p, "J", 0, []rune(""), false)
}

func TestMany1(t *testing.T) {
	p := Many1(Expect('H'))
	validate(t, "Many1(%q)", p, "HHHHH", 5, []rune("HHHHH"), false)
	validate(t, "Many1(%q)", p, "HHHHHJ", 5, []rune("HHHHH"), false)
	validate(t, "Many1(%q)", p, "H", 1, []rune("H"), false)
	validate(t, "Many1(%q)", p, "HJ", 1, []rune("H"), false)
	validate(t, "Many1(%q)", p, "", 0, []rune(""), true)
	validate(t, "Many1(%q)", p, "J", 0, []rune(""), true)
}

func TestManyN(t *testing.T) {
	p := ManyN(2, 5, Expect('H'))
	validate(t, "ManyN(%q)", p, "HHHHHHHJ", 5, []rune("HHHHH"), false)
	validate(t, "ManyN(%q)", p, "HHHHHHJ", 5, []rune("HHHHH"), false)
	validate(t, "ManyN(%q)", p, "HHHHHJ", 5, []rune("HHHHH"), false)
	validate(t, "ManyN(%q)", p, "HHHHJ", 4, []rune("HHHH"), false)
	validate(t, "ManyN(%q)", p, "HHHJ", 3, []rune("HHH"), false)
	validate(t, "ManyN(%q)", p, "HHJ", 2, []rune("HH"), false)
	validate(t, "ManyN(%q)", p, "HH", 2, []rune("HH"), false)
	validate(t, "ManyN(%q)", p, "HJ", 0, []rune(""), true)
	validate(t, "ManyN(%q)", p, "H", 0, []rune(""), true)
	validate(t, "ManyN(%q)", p, "J", 0, []rune(""), true)
	validate(t, "ManyN(%q)", p, "", 0, []rune(""), true)
}

func TestSeparatedList0(t *testing.T) {
	p := SeparatedList0(Expect(','), Expect('H'))
	validate(t, "SeparatedList0(%q)", p, "H,H,H", 5, []rune("HHH"), false)
	validate(t, "SeparatedList0(%q)", p, "H,H,", 3, []rune("HH"), false)
	validate(t, "SeparatedList0(%q)", p, "H,H", 3, []rune("HH"), false)
	validate(t, "SeparatedList0(%q)", p, "H,", 1, []rune("H"), false)
	validate(t, "SeparatedList0(%q)", p, "H", 1, []rune("H"), false)
	validate(t, "SeparatedList0(%q)", p, ",", 0, []rune(""), false)
	validate(t, "SeparatedList0(%q)", p, "H,H,J", 3, []rune("HH"), false)
	validate(t, "SeparatedList0(%q)", p, "", 0, []rune(""), false)
}

func TestSeparatedList1(t *testing.T) {
	p := SeparatedList1(Expect(','), Expect('H'))
	validate(t, "SeparatedList1(%q)", p, "H,H,H", 5, []rune("HHH"), false)
	validate(t, "SeparatedList1(%q)", p, "H,H,", 3, []rune("HH"), false)
	validate(t, "SeparatedList1(%q)", p, "H,H", 3, []rune("HH"), false)
	validate(t, "SeparatedList1(%q)", p, "H,", 1, []rune("H"), false)
	validate(t, "SeparatedList1(%q)", p, "H", 1, []rune("H"), false)
	validate(t, "SeparatedList1(%q)", p, ",", 0, []rune(""), true)
	validate(t, "SeparatedList1(%q)", p, "J,J", 0, []rune(""), true)
	validate(t, "SeparatedList1(%q)", p, "H,H,J", 3, []rune("HH"), false)
	validate(t, "SeparatedList1(%q)", p, "", 0, []rune(""), true)
}
