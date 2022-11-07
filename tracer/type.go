package tracer

import (
	"fmt"
	"sync"
)

type ArgType int

const (
	// public
	ArgTypeUnknown ArgType = iota
	ArgTypeData
	ArgTypeInt
	ArgTypeLong
	ArgTypeAddress
	ArgTypeUnsignedInt
	ArgTypeUnsignedLong
	ArgTypeObject
	ArgTypeErrorCode
	ArgTypeArray

	argStartInternal // TODO: use this to test all type converters are registered

	// internal
	argTypeString
	argTypeSockaddr
	argTypeIntOrErrorCode
	argTypeStat
	argTypePollFdArray
	argTypeSigAction
	argTypeIovecArray
	argTypeIntArray
	argTypeStringArray
	argTypeFdSet
	argTypeTimeval
	argTypeTimevalArray
	argTypeTimezone
	argTypeSHMIDDS
	argTypeTimespec
	argTypeTimespecArray
	argTypeItimerval
	argTypeMsghdr
	argTypeUnsignedIntPtr
	argTypeUnsignedInt64Ptr
	argTypeSockoptval
	argTypeWaitStatus
	argTypeRUsage
	argTypeRLimit
	argTypeUname
	argTypeSembuf
	argTypeSysinfo
	argTypeTms
	argTypeCapUserHeader
	argTypeCapUserData
	argTypeSigInfo
	argTypeStack
	argTypeUtimbuf
	argTypeUstat
	argTypeStatfs
	argTypeSchedParam
	argTypeUserDesc
	argTypeTimex
	argTypeIoEvent
	argTypeIoEvents
	argTypeIoCB
	argTypeItimerspec
	argTypeEpollEvent
	argTypeMqAttr
	argTypeMMsgHdrArray
	argTypeSchedAttr
	argTypeStatX
	argTypeIoUringParams
	argTypeCloneArgs
	argTypeOpenHow
	argTypeMountAttr
	argTypeLandlockRulesetAttr

	argEndInternal
)

type typeHandler func(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) error

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

func getHandler(t ArgType) typeHandler {
	typesRegistryMutex.RLock()
	defer typesRegistryMutex.RUnlock()
	return typesRegistry[t]
}

func handleType(arg *Arg, metadata ArgMetadata, raw, next, prev, ret uintptr, pid int) (err error) {
	typesRegistryMutex.RLock()
	defer typesRegistryMutex.RUnlock()
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error handling type %s with value 0x%x: %v", metadata.Name, raw, r)
		}
	}()
	if h, ok := typesRegistry[metadata.Type]; ok {
		return h(arg, metadata, raw, next, prev, ret, pid)
	}
	return nil
}
