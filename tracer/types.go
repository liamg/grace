package tracer

import (
	"fmt"
	"sync"
)

type ArgType int

const (
	ArgTypeUnknown ArgType = iota
	ArgTypeData
	ArgTypeInt
	ArgTypeStat
	ArgTypeLong
	ArgTypeAddress
	ArgTypeUnsignedInt
	ArgTypeUnsignedLong
	ArgTypePollFdArray
	ArgTypeObject
	ArgTypeErrorCode
	ArgTypeSigAction
	ArgTypeIovecArray
	ArgTypeIntArray
)

type typeHandler func(arg *Arg, metadata ArgMetadata, raw uintptr, next uintptr, ret uintptr, pid int) error

var typesRegistry = map[ArgType]typeHandler{}
var typesRegistryMutex = sync.RWMutex{}

func registerTypeHandler(t ArgType, h typeHandler) {
	typesRegistryMutex.Lock()
	defer typesRegistryMutex.Unlock()
	if _, ok := typesRegistry[t]; ok {
		panic(fmt.Sprintf("type handler for %d already registered", t))
	}
	typesRegistry[t] = h
}

func handleType(arg *Arg, metadata ArgMetadata, raw uintptr, next uintptr, ret uintptr, pid int) error {
	typesRegistryMutex.RLock()
	defer typesRegistryMutex.RUnlock()
	if h, ok := typesRegistry[metadata.Type]; ok {
		return h(arg, metadata, raw, next, ret, pid)
	}
	return nil
}
