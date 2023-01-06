package main

import (
	"fmt"
	"strconv"

	"github.com/jtdubs/go-nom"
	"github.com/jtdubs/go-nom/cache"
	"github.com/jtdubs/go-nom/printtracer"
	"github.com/jtdubs/go-nom/runes"
)

type Expr interface {
	Value() int
}

type BinaryExpr struct {
	L, R Expr
	Op   rune
}

func (b *BinaryExpr) String() string {
	return fmt.Sprintf("(%v %v %v)", b.L, string(b.Op), b.R)
}

func (b *BinaryExpr) Value() int {
	l, r := b.L.Value(), b.R.Value()
	switch b.Op {
	case '+':
		return l + r
	case '-':
		return l - r
	case '*':
		return l * r
	case '/':
		return l / r
	default:
		return 0
	}
}

type NumExpr struct {
	N int
}

func (n *NumExpr) String() string {
	return fmt.Sprint(n.N)
}

func (n *NumExpr) Value() int {
	return n.N
}

func CT[T any](p nom.ParseFn[rune, T]) nom.ParseFn[rune, T] {
	return nom.TraceN(1, cache.CacheN(1, p))
}

func Number(start nom.Cursor[rune]) (nom.Cursor[rune], Expr, error) {
	atoi := func(s string) Expr {
		n, _ := strconv.Atoi(s)
		return &NumExpr{n}
	}

	return CT(nom.Map(runes.Digit1, atoi))(start)
}

func Expression(start nom.Cursor[rune]) (nom.Cursor[rune], Expr, error) {
	return CT(nom.Alt(Parens, SumExpression))(start)
}

func Parens(start nom.Cursor[rune]) (nom.Cursor[rune], Expr, error) {
	return CT(runes.SurroundedBy('(', ')', Expression))(start)
}

func SumExpression(start nom.Cursor[rune]) (nom.Cursor[rune], Expr, error) {
	be := &BinaryExpr{}
	return CT(
		nom.Alt(
			nom.Value(Expr(be),
				runes.Phrase(
					nom.Bind(&be.L, ProductExpression),
					nom.Bind(&be.Op, SumOperator),
					nom.Bind(&be.R, SumExpression),
				),
			),
			ProductExpression,
		),
	)(start)
}

func SumOperator(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return CT(runes.OneOf("+-"))(start)
}

func ProductExpression(start nom.Cursor[rune]) (nom.Cursor[rune], Expr, error) {
	be := &BinaryExpr{}
	return CT(
		nom.Alt(
			nom.Value(Expr(be),
				runes.Phrase(
					nom.Bind(&be.L, Term),
					nom.Bind(&be.Op, ProductOperator),
					nom.Bind(&be.R, ProductExpression),
				),
			),
			Term,
		),
	)(start)
}

func ProductOperator(start nom.Cursor[rune]) (nom.Cursor[rune], rune, error) {
	return CT(runes.OneOf("*/"))(start)
}

func Term(start nom.Cursor[rune]) (nom.Cursor[rune], Expr, error) {
	return CT(nom.Alt(Number, Parens))(start)
}

func init() {
	nom.EnableTrace()
}

func main() {
	tracer := func() nom.Tracer[rune] {
		var opts printtracer.Options[rune]
		opts.IncludePackage("main")
		return opts.Tracer()
	}()

	start := runes.Cursor("    (1*7 + 1 + (2*3+	4/2))  ").WithTracer(tracer)
	rest, result, err := Expression(start)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Expression %v = %v\n", result, result.Value())
	if !rest.EOF() {
		fmt.Printf("Remaining: %q\n", string(rest.Rest()))
	}
}
