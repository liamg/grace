package filter

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/liamg/grace/tracer"
)

type Filter struct {
	allowNames   []string
	allowPaths   []string
	allowReturns []uint64
	failingOnly  bool
	passingOnly  bool
}

func Parse(input string) (*Filter, error) {
	filter := NewFilter()
	parts := strings.Split(input, "&")
	for _, part := range parts {
		if part == "" {
			continue
		}
		bits := strings.Split(part, "=")
		key := bits[0]
		value := bits[len(bits)-1]
		switch key {
		case "syscall", "name", "trace":
			filter.allowNames = append(filter.allowNames, strings.Split(value, ",")...)
		case "path":
			filter.allowPaths = append(filter.allowPaths, strings.Split(value, ",")...)
		case "ret", "retval", "return":
			ret, err := parseUint64(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse return value filter: %w", err)
			}
			filter.allowReturns = append(filter.allowReturns, ret)
		default:
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}
	}
	return filter, nil
}

func parseUint64(input string) (uint64, error) {
	if strings.HasPrefix(input, "0x") {
		return strconv.ParseUint(input[2:], 16, 64)
	}
	return strconv.ParseUint(input, 10, 64)
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) Match(call *tracer.Syscall, exit bool) bool {

	if len(f.allowNames) > 0 {
		var match bool
		for _, name := range f.allowNames {
			if name == call.Name() {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	if len(f.allowPaths) > 0 {
		var match bool
		for _, path := range f.allowPaths {
			for _, realPath := range call.Paths() {
				if realPath == path {
					match = true
					break
				}
			}
		}
		if !match {
			return false
		}
	}

	if len(f.allowReturns) > 0 {
		if !exit {
			return false
		}
		var match bool
		for _, ret := range f.allowReturns {
			if uintptr(ret) == call.Return().Raw() {
				match = true
				break
			}
		}
		if !match {
			return false
		}
	}

	if (f.passingOnly || f.failingOnly) && !exit {
		return false
	}

	if f.failingOnly && call.Return().Int() >= 0 {
		return false
	}

	if f.passingOnly && call.Return().Int() < 0 {
		return false
	}

	// TODO check more filters

	return true
}

func (f *Filter) SetFailingOnly(failing bool) {
	f.failingOnly = failing
}

func (f *Filter) SetPassingOnly(passing bool) {
	f.passingOnly = passing
}
