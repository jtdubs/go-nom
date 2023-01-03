package nom

type ParseFn[C comparable, T any] func(Cursor[C]) (Cursor[C], T, error)
