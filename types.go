package nom

import "context"

type ParseFn[C comparable, T any] func(context.Context, Cursor[C]) (Cursor[C], T, error)

type Tuple[T, U any] struct {
	A T
	B U
}
