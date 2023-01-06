package nom

import (
	"context"
	"errors"
)

func Any[C comparable](ctx context.Context, start Cursor[C]) (Cursor[C], C, error) {
	return Trace(func(ctx context.Context, start Cursor[C]) (Cursor[C], C, error) {
		if start.EOF() {
			return start, zero[C](), errors.New("EOF")
		}
		got := start.Read()
		return start.Next(), got, nil
	})(ctx, start)
}

func Rest[C comparable](ctx context.Context, start Cursor[C]) (Cursor[C], []C, error) {
	return Trace(func(ctx context.Context, start Cursor[C]) (Cursor[C], []C, error) {
		return start.ToEOF(), start.Rest(), nil
	})(ctx, start)
}

func RestLen[C comparable](ctx context.Context, start Cursor[C]) (Cursor[C], int, error) {
	return Trace(func(ctx context.Context, start Cursor[C]) (Cursor[C], int, error) {
		return start.ToEOF(), start.Len(), nil
	})(ctx, start)
}

func EOF[C comparable](ctx context.Context, start Cursor[C]) (Cursor[C], struct{}, error) {
	return Trace(func(ctx context.Context, start Cursor[C]) (Cursor[C], struct{}, error) {
		if !start.EOF() {
			return start, struct{}{}, errors.New("EOF() got slice, want EOF")
		}
		return start, struct{}{}, nil
	})(ctx, start)
}
