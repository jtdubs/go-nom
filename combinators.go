package nom

import (
	"errors"
	"fmt"
)

func zero[T any]() T {
	var result T
	return result
}

func Alt[C comparable, T any](ps ...ParseFn[C, T]) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		for _, p := range ps {
			end, result, err := p(start)
			if err != nil {
				continue
			}
			return end, result, nil
		}
		return start, zero[T](), errors.New("no alternatives matched")
	})
}

func Expect[C comparable](want C) ParseFn[C, C] {
	return Trace(func(start Cursor[C]) (Cursor[C], C, error) {
		if start.EOF() {
			return start, zero[C](), fmt.Errorf("got %v, want EOF", start.Read())
		}
		if got := start.Read(); got != want {
			return start, zero[C](), fmt.Errorf("got %v, want %v", got, want)
		}
		return start.Next(), want, nil
	})
}

func Expects[C comparable](want []C) ParseFn[C, []C] {
	fns := make([]ParseFn[C, C], len(want))
	for i, w := range want {
		fns[i] = Expect(w)
	}
	return Preceded(Seq(fns...), Success[C](want))
}

func Map[C comparable, T, U any](p ParseFn[C, T], fn func(T) U) ParseFn[C, U] {
	return Trace(func(start Cursor[C]) (Cursor[C], U, error) {
		end, result, err := p(start)
		if err != nil {
			return start, zero[U](), err
		}
		return end, fn(result), nil
	})
}

func Opt[C comparable, T any](p ParseFn[C, T]) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		end, res, err := p(start)
		if err != nil {
			return start, zero[T](), nil
		}
		return end, res, nil
	})
}

func Cond[C comparable, T any](b bool, p ParseFn[C, T]) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		if !b {
			return start, zero[T](), nil
		}
		return p(start)
	})
}

func Peek[C comparable, T any](p ParseFn[C, T]) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		_, res, err := p(start)
		return start, res, err
	})
}

func Verify[C comparable, T any](p ParseFn[C, T], checkFn func(T) bool) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		end, res, err := p(start)
		if err != nil {
			return start, zero[T](), err
		}
		if !checkFn(res) {
			return start, zero[T](), errors.New("Verify() check failed")
		}
		return end, res, nil
	})
}

func Value[C comparable, T, U any](p ParseFn[C, T], val U) ParseFn[C, U] {
	return Trace(func(start Cursor[C]) (Cursor[C], U, error) {
		end, _, err := p(start)
		if err != nil {
			return start, zero[U](), err
		}
		return end, val, nil
	})
}

func Not[C comparable, T any](p ParseFn[C, T]) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		_, _, err := p(start)
		if err != nil {
			return start, zero[T](), nil
		}
		return start, zero[T](), errors.New("Not()")
	})
}

func Success[C comparable, T any](val T) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		return start, val, nil
	})
}

func Failure[C comparable, T any](err error) ParseFn[C, T] {
	return Trace(func(start Cursor[C]) (Cursor[C], T, error) {
		return start, zero[T](), err
	})
}

func Recognize[C comparable, T any](p ParseFn[C, T]) ParseFn[C, []C] {
	return Trace(func(start Cursor[C]) (Cursor[C], []C, error) {
		end, _, err := p(start)
		if err != nil {
			return start, nil, err
		}
		return end, start.To(end), nil
	})
}

func Bind[C comparable, T any](p ParseFn[C, T], place *T) ParseFn[C, struct{}] {
	return Trace(func(start Cursor[C]) (Cursor[C], struct{}, error) {
		end, res, err := p(start)
		if err != nil {
			return start, struct{}{}, err
		}
		*place = res
		return end, struct{}{}, err
	})
}

func Discard[C comparable, T any](p ParseFn[C, T]) ParseFn[C, struct{}] {
	return Trace(func(start Cursor[C]) (Cursor[C], struct{}, error) {
		end, _, err := p(start)
		return end, struct{}{}, err
	})
}

func Satisfy[C comparable](testFn func(C) bool) ParseFn[C, C] {
	return Trace(func(start Cursor[C]) (Cursor[C], C, error) {
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
