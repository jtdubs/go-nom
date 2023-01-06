package nom

type ParseFn[C comparable, T any] func(Cursor[C]) (Cursor[C], T, error)

type Tuple[T, U any] struct {
	A T
	B U
}
