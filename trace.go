package nom

import (
	"context"
	"runtime"
	"strings"
)

type ContextKeyType int

const (
	TracerKey ContextKeyType = iota
	TraceEnabledKey
)

func WithTracer[T comparable](ctx context.Context, tracer Tracer[T]) context.Context {
	return context.WithValue(ctx, TracerKey, tracer)
}

func WithTracing(ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceEnabledKey, true)
}

func WithoutTracing(ctx context.Context) context.Context {
	return context.WithValue(ctx, TraceEnabledKey, false)
}

type Tracer[T comparable] interface {
	Enter(ctx context.Context, name string, start Cursor[T])
	Exit(ctx context.Context, name string, start, end Cursor[T], result any, err error)
}

var traceSupported = false

func TraceSupported() {
	traceSupported = true
}

func Trace[C comparable, T any](fn ParseFn[C, T]) ParseFn[C, T] {
	return TraceN(1, fn)
}

func TraceN[C comparable, T any](depth int, fn ParseFn[C, T]) ParseFn[C, T] {
	if !traceSupported {
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

	return func(ctx context.Context, start Cursor[C]) (end Cursor[C], res T, err error) {
		tracer, ok := ctx.Value(TracerKey).(Tracer[C])
		tracingEnabled, _ := ctx.Value(TraceEnabledKey).(bool)
		if ok && tracingEnabled {
			tracer.Enter(ctx, name, start)
		}
		end, res, err = fn(ctx, start)
		if ok && tracingEnabled {
			tracer.Exit(ctx, name, start, end, res, err)
		}
		return
	}
}
