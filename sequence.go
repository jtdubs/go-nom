package nom

func Seq[C comparable, T any](ps ...ParseFn[C, T]) ParseFn[C, []T] {
	return Trace(func(start Cursor[C]) (Cursor[C], []T, error) {
		var results []T
		end := start
		for _, p := range ps {
			var (
				result T
				err    error
			)
			end, result, err = p(end)
			if err != nil {
				return start, nil, err
			}
			results = append(results, result)
		}
		return end, results, nil
	})
}

func Surrounded[C comparable, F, L, M any](first ParseFn[C, F], last ParseFn[C, L], middle ParseFn[C, M]) ParseFn[C, M] {
	return Trace(func(start Cursor[C]) (Cursor[C], M, error) {
		var (
			res M
			err error
		)
		end := start
		if end, _, err = first(end); err != nil {
			return start, zero[M](), err
		}
		if end, res, err = middle(end); err != nil {
			return start, zero[M](), err
		}
		if end, _, err = last(end); err != nil {
			return start, zero[M](), err
		}
		return end, res, nil
	})
}

func Preceded[C comparable, A, B any](first ParseFn[C, A], second ParseFn[C, B]) ParseFn[C, B] {
	return Trace(func(start Cursor[C]) (Cursor[C], B, error) {
		var (
			res B
			err error
		)
		end := start
		if end, _, err = first(end); err != nil {
			return start, zero[B](), err
		}
		if end, res, err = second(end); err != nil {
			return start, zero[B](), err
		}
		return end, res, nil
	})
}

func Terminated[C comparable, A, B any](first ParseFn[C, A], second ParseFn[C, B]) ParseFn[C, A] {
	return Trace(func(start Cursor[C]) (Cursor[C], A, error) {
		var (
			res A
			err error
		)
		end := start
		end, res, err = first(end)
		if err != nil {
			return start, zero[A](), err
		}
		end, _, err = second(end)
		if err != nil {
			return start, zero[A](), err
		}
		return end, res, nil
	})
}
