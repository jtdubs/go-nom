package nom

import (
	"runtime"
	"strings"
)

var (
	traceEnabled bool
)

type Tracer[T comparable] interface {
	Enter(name string, cursor Cursor[T])
	Exit(name string, oldCursor, newCursor Cursor[T], result any, err error)
}

func EnableTrace() {
	traceEnabled = true
}

func DisableTrace() {
	traceEnabled = false
}

func Trace[C comparable, T any](fn ParseFn[C, T]) ParseFn[C, T] {
	return TraceN(1, fn)
}

func TraceN[C comparable, T any](depth int, fn ParseFn[C, T]) ParseFn[C, T] {
	if !traceEnabled {
		return fn
	}

	pc, _, _, ok := runtime.Caller(depth + 1)
	parent := runtime.FuncForPC(pc)
	name := "unknown"
	if ok && parent != nil {
		name = parent.Name()
		if idx := strings.IndexRune(name, '['); idx != -1 {
			name = name[:idx]
		}
		if idx := strings.LastIndex(name, "/"); idx != -1 {
			name = name[idx+1:]
		}
	}

	return func(start Cursor[C]) (end Cursor[C], res T, err error) {
		tracer := start.Tracer()
		if tracer != nil {
			tracer.Enter(name, start)
		}
		end, res, err = fn(start)
		if tracer != nil {
			tracer.Exit(name, start, end, res, err)
		}
		return
	}
}
