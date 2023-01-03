package nom

import "errors"

func Any[C comparable](start Cursor[C]) (Cursor[C], C, error) {
	return Trace(func(start Cursor[C]) (Cursor[C], C, error) {
		if start.EOF() {
			return start, zero[C](), errors.New("EOF")
		}
		got := start.Read()
		return start.Next(), got, nil
	})(start)
}

func Rest[C comparable](start Cursor[C]) (Cursor[C], []C, error) {
	return Trace(func(start Cursor[C]) (Cursor[C], []C, error) {
		return start.ToEOF(), start.Rest(), nil
	})(start)
}

func RestLen[C comparable](start Cursor[C]) (Cursor[C], int, error) {
	return Trace(func(start Cursor[C]) (Cursor[C], int, error) {
		return start.ToEOF(), start.Len(), nil
	})(start)
}

func EOF[C comparable](start Cursor[C]) (Cursor[C], struct{}, error) {
	return Trace(func(start Cursor[C]) (Cursor[C], struct{}, error) {
		if !start.EOF() {
			return start, struct{}{}, errors.New("EOF() got slice, want EOF")
		}
		return start, struct{}{}, nil
	})(start)
}
