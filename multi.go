package nom

import (
	"fmt"
)

func Many0[C comparable, T any](p ParseFn[C, T]) ParseFn[C, []T] {
	return Trace(func(start Cursor[C]) (end Cursor[C], results []T, err error) {
		end = start
		for {
			var res T
			end, res, err = p(end)
			if err != nil {
				return end, results, nil
			}
			results = append(results, res)
		}
	})
}

func Many1[C comparable, T any](p ParseFn[C, T]) ParseFn[C, []T] {
	return Trace(func(start Cursor[C]) (Cursor[C], []T, error) {
		end, res, err := p(start)
		if err != nil {
			return start, nil, err
		}
		var results []T
		for err == nil {
			results = append(results, res)
			end, res, err = p(end)
		}
		return end, results, nil
	})
}

func ManyN[C comparable, T any](min, max int, p ParseFn[C, T]) ParseFn[C, []T] {
	return Trace(func(start Cursor[C]) (Cursor[C], []T, error) {
		end := start
		var (
			results []T
			err     error
		)
		for len(results) < max {
			var res T
			end, res, err = p(end)
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

func SeparatedList0[C comparable, T, D any](delim ParseFn[C, D], values ParseFn[C, T]) ParseFn[C, []T] {
	return Trace(func(start Cursor[C]) (Cursor[C], []T, error) {
		var results []T
		end, res, err := values(start)
		if err != nil {
			return start, nil, nil
		}
		results = append(results, res)
		for {
			delimEnd, _, err := delim(end)
			if err != nil {
				return end, results, nil
			}
			valueEnd, res, err := values(delimEnd)
			if err != nil {
				return end, results, nil
			}
			end = valueEnd
			results = append(results, res)
		}
	})
}

func SeparatedList1[C comparable, T, D any](delim ParseFn[C, D], values ParseFn[C, T]) ParseFn[C, []T] {
	return Trace(func(start Cursor[C]) (Cursor[C], []T, error) {
		end, res, err := values(start)
		if err != nil {
			return start, nil, err
		}
		results := []T{res}
		for {
			delimEnd, _, err := delim(end)
			if err != nil {
				return end, results, nil
			}
			valueEnd, res, err := values(delimEnd)
			if err != nil {
				return end, results, nil
			}
			end = valueEnd
			results = append(results, res)
		}
	})
}
