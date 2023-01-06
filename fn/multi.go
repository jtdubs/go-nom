package fn

import (
	"context"
	"fmt"

	"github.com/jtdubs/go-nom"
	"github.com/jtdubs/go-nom/trace"
)

func Many0[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, []T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (end nom.Cursor[C], results []T, err error) {
		end = start
		for {
			var res T
			end, res, err = p(ctx, end)
			if err != nil {
				return end, results, nil
			}
			results = append(results, res)
		}
	})
}

func Many1[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, []T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []T, error) {
		end, res, err := p(ctx, start)
		if err != nil {
			return start, nil, err
		}
		var results []T
		for err == nil {
			results = append(results, res)
			end, res, err = p(ctx, end)
		}
		return end, results, nil
	})
}

func ManyN[C comparable, T any](min, max int, p nom.ParseFn[C, T]) nom.ParseFn[C, []T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []T, error) {
		end := start
		var (
			results []T
			err     error
		)
		for len(results) < max {
			var res T
			end, res, err = p(ctx, end)
			if err != nil {
				break
			}
			results = append(results, res)
		}
		if len(results) < min {
			return start, nil, fmt.Errorf("ManyN() got %v, wanted [%v, %v]", len(results), min, max)
		}
		return end, results, nil
	})
}

func ManyTill[C comparable, T, U any](f nom.ParseFn[C, T], g nom.ParseFn[C, U]) nom.ParseFn[C, nom.Tuple[[]T, U]] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (end nom.Cursor[C], res nom.Tuple[[]T, U], err error) {
		end = start
		for {
			var (
				u U
				t T
			)
			if end, u, err = g(ctx, end); err == nil {
				res.B = u
				return
			}
			if end, t, err = f(ctx, end); err != nil {
				end = start
				res.A = nil
				return
			}
			res.A = append(res.A, t)
		}
	})
}

func SeparatedList0[C comparable, T, D any](delim nom.ParseFn[C, D], values nom.ParseFn[C, T]) nom.ParseFn[C, []T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []T, error) {
		var results []T
		end, res, err := values(ctx, start)
		if err != nil {
			return start, nil, nil
		}
		results = append(results, res)
		for {
			delimEnd, _, err := delim(ctx, end)
			if err != nil {
				return end, results, nil
			}
			valueEnd, res, err := values(ctx, delimEnd)
			if err != nil {
				return end, results, nil
			}
			end = valueEnd
			results = append(results, res)
		}
	})
}

func SeparatedList1[C comparable, T, D any](delim nom.ParseFn[C, D], values nom.ParseFn[C, T]) nom.ParseFn[C, []T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []T, error) {
		end, res, err := values(ctx, start)
		if err != nil {
			return start, nil, err
		}
		results := []T{res}
		for {
			delimEnd, _, err := delim(ctx, end)
			if err != nil {
				return end, results, nil
			}
			valueEnd, res, err := values(ctx, delimEnd)
			if err != nil {
				return end, results, nil
			}
			end = valueEnd
			results = append(results, res)
		}
	})
}
