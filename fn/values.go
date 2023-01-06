package fn

import (
	"context"
	"errors"

	"github.com/jtdubs/go-nom"
	"github.com/jtdubs/go-nom/trace"
)

func Any[C comparable](ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], C, error) {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], C, error) {
		if start.EOF() {
			return start, zero[C](), errors.New("EOF")
		}
		got := start.Read()
		return start.Next(), got, nil
	})(ctx, start)
}

func Rest[C comparable](ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []C, error) {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], []C, error) {
		return start.ToEOF(), start.Rest(), nil
	})(ctx, start)
}

func RestLen[C comparable](ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], int, error) {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], int, error) {
		return start.ToEOF(), start.Len(), nil
	})(ctx, start)
}

func EOF[C comparable](ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], struct{}, error) {
	return trace.Trace(func(ctx context.Context, start nom.Cursor[C]) (nom.Cursor[C], struct{}, error) {
		if !start.EOF() {
			return start, struct{}{}, errors.New("EOF() got slice, want EOF")
		}
		return start, struct{}{}, nil
	})(ctx, start)
}
