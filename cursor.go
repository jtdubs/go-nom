package nom

import (
	"fmt"
)

type Cursor[T comparable] struct {
	buffer []T
	offset int
	tracer Tracer[T]
}

type Span[T comparable] struct {
	Start, End Cursor[T]
}

func NewCursor[T comparable](ts []T) Cursor[T] {
	return Cursor[T]{
		buffer: ts,
		offset: 0,
		tracer: nil,
	}
}

func (c Cursor[T]) String() string {
	return fmt.Sprintf("cursor(%v)", c.offset)
}

func (c Cursor[T]) WithTracer(t Tracer[T]) Cursor[T] {
	return Cursor[T]{
		buffer: c.buffer,
		offset: c.offset,
		tracer: t,
	}
}

func (c Cursor[T]) Tracer() Tracer[T] {
	return c.tracer
}

func (c Cursor[T]) EOF() bool {
	return c.offset >= len(c.buffer)
}

func (c Cursor[T]) Rest() []T {
	if c.EOF() {
		return nil
	}
	return c.buffer[c.offset:]
}

func (c Cursor[T]) Len() int {
	return len(c.buffer) - c.offset
}

func (c Cursor[T]) Position() int {
	return c.offset
}

func (c Cursor[T]) Read() (value T) {
	if c.offset >= len(c.buffer) {
		return
	}
	value = c.buffer[c.offset]
	return
}

func (c Cursor[T]) Next() Cursor[T] {
	if c.EOF() {
		return c
	}
	return Cursor[T]{
		buffer: c.buffer,
		offset: c.offset + 1,
		tracer: c.tracer,
	}
}

func (c Cursor[T]) ToEOF() Cursor[T] {
	return Cursor[T]{
		buffer: c.buffer,
		offset: len(c.buffer),
		tracer: c.tracer,
	}
}

func (c Cursor[T]) To(other Cursor[T]) []T {
	if c.EOF() || &c.buffer[0] != &other.buffer[0] {
		return nil
	}
	return c.buffer[c.offset:other.offset]
}

func (c Cursor[T]) Addr() *T {
	if c.EOF() {
		return nil
	}
	return &c.buffer[c.offset]
}
