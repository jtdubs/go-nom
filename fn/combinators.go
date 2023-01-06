package fn

import (
	"context"
	"errors"
	"fmt"

	"github.com/jtdubs/go-nom"
	"github.com/jtdubs/go-nom/trace"
)

func zero[T any]() T {
	var result T
	return result
}

func Alt[C comparable, T any](ps ...nom.ParseFn[C, T]) nom.ParseFn[C, T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		for _, p := range ps {
			end, result, err := p(ctx, start)
			if err != nil {
				continue
			}
			return end, result, nil
		}
		return start, zero[T](), errors.New("no alternatives matched")
	})
}

func Expect[C comparable](want C) nom.ParseFn[C, C] {
	return trace.Trace(func(_ context.Context, start nom.Cursor[C]) (nom.Cursor[C], C, error) {
		if start.EOF() {
			return start, zero[C](), fmt.Errorf("got %v, want EOF", start.Read())
		}
		if got := start.Read(); got != want {
			return start, zero[C](), fmt.Errorf("got %v, want %v", got, want)
		}
		return start.Next(), want, nil
	})
}

func Expects[C comparable](want []C) nom.ParseFn[C, []C] {
	fns := make([]nom.ParseFn[C, C], len(want))
	for i, w := range want {
		fns[i] = Expect(w)
	}
	return trace.Trace(Seq(fns...))
}

func Map[C comparable, T, U any](p nom.ParseFn[C, T], fn func(T) U) nom.ParseFn[C, U] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], U, error) {
		end, result, err := p(ctx, start)
		if err != nil {
			return start, zero[U](), err
		}
		return end, fn(result), nil
	})
}

func Opt[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		end, res, err := p(ctx, start)
		if err != nil {
			return start, zero[T](), nil
		}
		return end, res, nil
	})
}

func Cond[C comparable, T any](b bool, p nom.ParseFn[C, T]) nom.ParseFn[C, T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		if !b {
			return start, zero[T](), nil
		}
		return p(ctx, start)
	})
}

func Peek[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		_, res, err := p(ctx, start)
		return start, res, err
	})
}

func Verify[C comparable, T any](p nom.ParseFn[C, T], checkFn func(T) bool) nom.ParseFn[C, T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		end, res, err := p(ctx, start)
		if err != nil {
			return start, zero[T](), err
		}
		if !checkFn(res) {
			return start, zero[T](), errors.New("Verify() check failed")
		}
		return end, res, nil
	})
}

func Value[C comparable, T, U any](val U, p nom.ParseFn[C, T]) nom.ParseFn[C, U] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], U, error) {
		end, _, err := p(ctx, start)
		if err != nil {
			return start, zero[U](), err
		}
		return end, val, nil
	})
}

func Not[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, T] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		_, _, err := p(ctx, start)
		if err != nil {
			return start, zero[T](), nil
		}
		return start, zero[T](), errors.New("Not()")
	})
}

func Success[C comparable, T any](val T) nom.ParseFn[C, T] {
	return trace.Trace(func(_ context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		return start, val, nil
	})
}

func Failure[C comparable, T any](err error) nom.ParseFn[C, T] {
	return trace.Trace(func(_ context.Context, start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		return start, zero[T](), err
	})
}

func Recognize[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, []C] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []C, error) {
		end, _, err := p(ctx, start)
		if err != nil {
			return start, nil, err
		}
		return end, start.To(end), nil
	})
}

func Bind[C comparable, T any](place *T, p nom.ParseFn[C, T]) nom.ParseFn[C, struct{}] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], struct{}, error) {
		end, res, err := p(ctx, start)
		if err != nil {
			return start, struct{}{}, err
		}
		*place = res
		return end, struct{}{}, err
	})
}

func Discard[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, struct{}] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], struct{}, error) {
		end, _, err := p(ctx, start)
		return end, struct{}{}, err
	})
}

func Satisfy[C comparable](testFn func(C) bool) nom.ParseFn[C, C] {
	return trace.Trace(func(_ context.Context, start nom.Cursor[C]) (nom.Cursor[C], C, error) {
		if start.EOF() {
			return start, zero[C](), errors.New("EOF")
		}
		got := start.Read()
		if !testFn(got) {
			return start, zero[C](), fmt.Errorf("%v does not satisfy test", got)
		}
		return start.Next(), got, nil
	})
}

func Pair[C comparable, T, U any](t nom.ParseFn[C, T], u nom.ParseFn[C, U]) nom.ParseFn[C, nom.Tuple[T, U]] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], nom.Tuple[T, U], error) {
		end, a, err := t(ctx, start)
		if err != nil {
			return start, zero[nom.Tuple[T, U]](), err
		}
		end, b, err := u(ctx, end)
		if err != nil {
			return start, zero[nom.Tuple[T, U]](), err
		}
		return end, nom.Tuple[T, U]{a, b}, nil
	})
}

func First[C comparable, T, U any](p nom.ParseFn[C, nom.Tuple[T, U]]) nom.ParseFn[C, T] {
	return trace.Trace(Map(p, func(t nom.Tuple[T, U]) T { return t.A }))
}

func Second[C comparable, T, U any](p nom.ParseFn[C, nom.Tuple[T, U]]) nom.ParseFn[C, U] {
	return trace.Trace(Map(p, func(t nom.Tuple[T, U]) U { return t.B }))
}

func Spanning[C comparable, T any](p nom.ParseFn[C, T]) nom.ParseFn[C, nom.Span[C]] {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (end nom.Cursor[C], res nom.Span[C], err error) {
		res.Start = start
		end, _, err = p(ctx, start)
		res.End = end
		return
	})
}
