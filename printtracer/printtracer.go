package printtracer

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/jtdubs/go-nom"
)

type Options[T comparable] struct {
	allowlist map[*regexp.Regexp]struct{}
	blocklist map[*regexp.Regexp]struct{}
}

func (pt *Options[T]) IncludePackage(packages ...string) {
	for _, p := range packages {
		pt.Include(fmt.Sprintf("%v\\..*", p))
	}
}

func (pt *Options[T]) ExcludePackage(packages ...string) {
	for _, p := range packages {
		pt.Exclude(fmt.Sprintf("%v\\..*", p))
	}
}

func (pt *Options[T]) Include(patterns ...string) {
	if pt.allowlist == nil {
		pt.allowlist = make(map[*regexp.Regexp]struct{})
	}
	for _, p := range patterns {
		pt.allowlist[regexp.MustCompile(p)] = struct{}{}
	}
}

func (pt *Options[T]) Exclude(patterns ...string) {
	if pt.blocklist == nil {
		pt.blocklist = make(map[*regexp.Regexp]struct{})
	}
	for _, p := range patterns {
		pt.blocklist[regexp.MustCompile(p)] = struct{}{}
	}
}

func (pt *Options[T]) Tracer() nom.Tracer[T] {
	if pt.allowlist == nil {
		pt.allowlist = make(map[*regexp.Regexp]struct{})
	}
	if pt.blocklist == nil {
		pt.blocklist = make(map[*regexp.Regexp]struct{})
	}
	return &tracer[T]{
		level:     0,
		allowlist: pt.allowlist,
		blocklist: pt.blocklist,
	}
}

type tracer[T comparable] struct {
	level     int
	allowlist map[*regexp.Regexp]struct{}
	blocklist map[*regexp.Regexp]struct{}
}

func New[T comparable]() nom.Tracer[T] {
	return &tracer[T]{0, make(map[*regexp.Regexp]struct{}), make(map[*regexp.Regexp]struct{})}
}

func (pt *tracer[T]) skip(name string) bool {
	for p := range pt.blocklist {
		if p.MatchString(name) {
			return true
		}
	}

	if len(pt.allowlist) == 0 {
		return false
	}

	for p := range pt.allowlist {
		if p.MatchString(name) {
			return false
		}
	}
	return true
}

func (bt *tracer[T]) Enter(_ context.Context, name string, start nom.Cursor[T]) {
	if bt.skip(name) {
		return
	}
	for i := 0; i < bt.level; i++ {
		fmt.Print("  ")
	}
	fmt.Printf("%v(%v)\n", name, start.Position())
	bt.level = bt.level + 1
}

func (bt *tracer[T]) Exit(_ context.Context, name string, start, end nom.Cursor[T], result any, err error) {
	if bt.skip(name) {
		return
	}
	bt.level = bt.level - 1
	for i := 0; i < bt.level; i++ {
		fmt.Print("  ")
	}
	if err == nil {
		switch result.(type) {
		case rune, string:
			fmt.Printf("< %q", result)
		default:
			fmt.Printf("< %v", result)
		}
		span := reflect.ValueOf(start.To(end)).Interface()
		switch rs := span.(type) {
		case []rune:
			s := string(rs)
			if len(s) > 20 {
				fmt.Printf(" [from %q]", string(s[:10])+"..."+string(s[len(s)-10:]))
			} else {
				fmt.Printf(" [from %q]", string(s))
			}
		default:
		}
		fmt.Printf("\n")
	} else {
		fmt.Printf("! %v\n", err)
	}
}
