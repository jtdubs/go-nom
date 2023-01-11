package fn

import (
	"context"

	"github.com/jtdubs/go-nom"
	"github.com/jtdubs/go-nom/trace"
)

func Seq[C comparable, T any](ps ...nom.ParseFn[C, T]) nom.ParseFn[C, []T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []T, error) {
		var results []T
		end := start
		for _, p := range ps {
			var (
				result T
				err    error
			)
			end, result, err = p(ctx, end)
			if err != nil {
				return start, nil, err
			}
			results = append(results, result)
		}
		return end, results, nil
	})
}

func Surrounded[C comparable, F, L, M any](first nom.ParseFn[C, F], last nom.ParseFn[C, L], middle nom.ParseFn[C, M]) nom.ParseFn[C, M] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], M, error) {
		var (
			res M
			err error
		)
		end := start
		if end, _, err = first(ctx, end); err != nil {
			return start, zero[M](), err
		}
		if end, res, err = middle(ctx, end); err != nil {
			return start, zero[M](), err
		}
		if end, _, err = last(ctx, end); err != nil {
			return start, zero[M](), err
		}
		return end, res, nil
	})
}

func Preceded[C comparable, A, B any](first nom.ParseFn[C, A], second nom.ParseFn[C, B]) nom.ParseFn[C, B] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], B, error) {
		var (
			res B
			err error
		)
		end := start
		if end, _, err = first(ctx, end); err != nil {
			return start, zero[B](), err
		}
		if end, res, err = second(ctx, end); err != nil {
			return start, zero[B](), err
		}
		return end, res, nil
	})
}

func Terminated[C comparable, A, B any](first nom.ParseFn[C, A], second nom.ParseFn[C, B]) nom.ParseFn[C, A] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], A, error) {
		var (
			res A
			err error
		)
		end := start
		end, res, err = first(ctx, end)
		if err != nil {
			return start, zero[A](), err
		}
		end, _, err = second(ctx, end)
		if err != nil {
			return start, zero[A](), err
		}
		return end, res, nil
	})
}

func Append[C comparable, A any](first nom.ParseFn[C, []A], rest ...nom.ParseFn[C, A]) nom.ParseFn[C, []A] {
	return trace.Trace(Map(Pair(first, Seq(rest...)), func(t nom.Tuple[[]A, []A]) []A { return append(t.A, t.B...) }))
}
