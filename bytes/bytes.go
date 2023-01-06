package bytes

import (
	"context"
	"errors"
	"fmt"

	"github.com/jtdubs/go-nom"
)

func Byte(want byte) nom.ParseFn[byte, byte] {
	return func(_ context.Context, start nom.Cursor[byte]) (nom.Cursor[byte], byte, error) {
		if start.EOF() {
			return start, 0, fmt.Errorf("got %v, want EOF", start.Read())
		}
		if got := start.Read(); got != want {
			return start, 0, fmt.Errorf("got %v, want %v", got, want)
		}
		return start.Next(), want, nil
	}
}

func Tag(tag string) nom.ParseFn[byte, string] {
	var bytes []nom.ParseFn[byte, byte]
	for _, b := range []byte(tag) {
		bytes = append(bytes, Byte(b))
	}
	return nom.Preceded(nom.Seq(bytes...), nom.Success[byte](tag))
}

func Satisfy(testFn func(byte) bool) nom.ParseFn[byte, byte] {
	return func(_ context.Context, start nom.Cursor[byte]) (nom.Cursor[byte], byte, error) {
		if start.EOF() {
			return start, 0, errors.New("EOF")
		}
		got := start.Read()
		if !testFn(got) {
			return start, 0, fmt.Errorf("%v does not satisfy test", got)
		}
		return start.Next(), got, nil
	}
}

func OneOf(allowlist []byte) nom.ParseFn[byte, byte] {
	lookup := map[byte]struct{}{}
	for _, b := range allowlist {
		lookup[b] = struct{}{}
	}

	return func(_ context.Context, start nom.Cursor[byte]) (nom.Cursor[byte], byte, error) {
		if start.EOF() {
			return start, 0, errors.New("EOF")
		}
		got := start.Read()
		if _, ok := lookup[got]; !ok {
			return start, 0, fmt.Errorf("%q not allowed", got)
		}
		return start.Next(), got, nil
	}
}

func NoneOf(blocklist []byte) nom.ParseFn[byte, byte] {
	lookup := map[byte]struct{}{}
	for _, b := range blocklist {
		lookup[b] = struct{}{}
	}

	return func(_ context.Context, start nom.Cursor[byte]) (nom.Cursor[byte], byte, error) {
		if start.EOF() {
			return start, 0, errors.New("EOF")
		}
		got := start.Read()
		if _, ok := lookup[got]; ok {
			return start, 0, fmt.Errorf("%q not allowed", got)
		}
		return start.Next(), got, nil
	}
}

func Any() nom.ParseFn[byte, byte] {
	return func(_ context.Context, start nom.Cursor[byte]) (nom.Cursor[byte], byte, error) {
		if start.EOF() {
			return start, 0, errors.New("EOF")
		}
		got := start.Read()
		return start.Next(), got, nil
	}
}

func Cursor(b []byte) nom.Cursor[byte] {
	return nom.NewCursor(b)
}
