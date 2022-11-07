package tracer

import (
	"unsafe"

	"github.com/liamg/grace/tracer/annotation"
)

type ArgMetadata struct {
	Name        string
	Type        ArgType
	Annotator   func(arg annotation.Arg, pid int)
	LenSource   LenSource
	Optional    bool
	Destination bool
	FixedCount  int
}

type LenSource uint8

const (
	LenSourceNone LenSource = iota
	LenSourcePrev
	LenSourceNext
	LenSourceNextPointer
	LenSourceReturnValue
	LenSourceFixed
)

type ReturnMetadata ArgMetadata

type Arg struct {
	name       string
	t          ArgType
	raw        uintptr
	data       []byte
	annotation string
	replace    bool // replace value output with annotation
	bitSize    int
	obj        *Object
	array      []Arg
	known      bool
}

type Object struct {
	Name       string
	Properties []Arg
}

func (s Arg) Known() bool {
	return s.known
}

func (s Arg) Name() string {
	return s.name
}

func (s Arg) Type() ArgType {
	return s.t
}

func (s Arg) Raw() uintptr {
	return s.raw
}

func (s Arg) Int() int {
	switch s.t {
	case ArgTypeLong:
		switch s.bitSize {
		case 32:
			return int(int32(s.raw))
		case 64:
			return int(int64(s.raw))
		}
	case ArgTypeInt, ArgTypeUnknown:
		// an int is 32-bit on both 32-bit and 64-bit linux
		return int(int32(s.raw))
	}
	return int(s.raw)
}

func (s Arg) Data() []byte {
	return s.data
}

func (s Arg) Annotation() string {
	return s.annotation
}

func (s Arg) ReplaceValueWithAnnotation() bool {
	return s.replace
}

func (s Arg) Object() *Object {
	return s.obj
}

func (s Arg) Array() []Arg {
	return s.array
}

func (s *Arg) SetAnnotation(annotation string, replace bool) {
	s.annotation = annotation
	s.replace = replace
}

func processArgument(raw uintptr, next, prev uintptr, ret uintptr, metadata ArgMetadata, pid int, exit bool) (*Arg, error) {
	arg := &Arg{
		name:    metadata.Name,
		t:       metadata.Type,
		raw:     raw,
		bitSize: bitSize,
		known:   true,
	}

	// if we're on the syscall enter and the argument is a pointer for a destination, we don't know the value yet
	if !exit && metadata.Destination {
		arg.known = false
		return arg, nil
	}

	// resolve next to int from next pointer
	if metadata.LenSource == LenSourceNextPointer && (exit || !metadata.Destination) && next > 0 {
		var realNext uint32
		buf, err := readSize(pid, next, unsafe.Sizeof(realNext))
		if err != nil {
			return nil, err
		}
		next = uintptr(decodeInt(buf))
	}

	// process the argument data into something meaningful
	if err := handleType(arg, metadata, raw, next, prev, ret, pid); err != nil {
		return nil, err
	}

	// always apply annotations
	if metadata.Annotator != nil {
		metadata.Annotator(arg, pid)
	}

	return arg, nil
}
