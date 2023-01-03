package cache

import (
	"log"
	"runtime"

	"github.com/jtdubs/go-nom"
)

type cacheValue[C comparable, T any] struct {
	end   nom.Cursor[C]
	value T
	err   error
}

var caches = make(map[string]any)

func Cache[C comparable, T any](fn nom.ParseFn[C, T]) nom.ParseFn[C, T] {
	return CacheN(1, fn)
}

func CacheN[C comparable, T any](skip int, fn nom.ParseFn[C, T]) nom.ParseFn[C, T] {
	// Get caller name
	pc, file, line, ok := runtime.Caller(skip + 1)
	parent := runtime.FuncForPC(pc)
	if !ok && parent == nil {
		log.Printf("Cache failed: unable to determine function for %v:%v", file, line)
		return fn
	}

	// Create cache
	name := parent.Name()
	if _, ok := caches[name]; !ok {
		caches[name] = map[*C]cacheValue[C, T]{}
	}
	cache := caches[name].(map[*C]cacheValue[C, T])
	if cache == nil {
		log.Printf("Cache failed: incompatible cache found for %q", name)
		return fn
	}

	return func(start nom.Cursor[C]) (nom.Cursor[C], T, error) {
		cacheVal, ok := cache[start.Addr()]
		if !ok {
			end, res, err := fn(start)
			cacheVal = cacheValue[C, T]{end, res, err}
			cache[start.Addr()] = cacheVal
		}
		return cacheVal.end, cacheVal.value, cacheVal.err
	}
}
