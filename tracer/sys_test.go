package tracer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SyscallSupport(t *testing.T) {

	for number, meta := range sysMap {
		t.Run(fmt.Sprintf("syscall %d: %s", number, meta.Name), func(t *testing.T) {
			checkSyscall(t, number, meta)
		})
	}
}

func checkSyscall(t *testing.T, number int, meta SyscallMetadata) {
	require.NotNil(t, meta)
	assert.NotEmpty(t, meta.Name)
	assert.NotEqualf(t, ArgTypeUnknown, meta.ReturnValue.Type, "syscall %d (%s) has unspecified return value type", number, meta.Name)
	for i, arg := range meta.Args {
		assert.NotEqualf(t, ArgTypeUnknown, arg.Type, "syscall %d (%s) has unspecified argument type", number, meta.Name)
		switch arg.Type {
		case ArgTypeData:
			assert.NotEqual(t, LenSourceNone, arg.LenSource)
		}
		switch arg.LenSource {
		case LenSourceFixed:
			assert.NotZero(t, arg.FixedCount)
		case LenSourceNextPointer, LenSourceNext:
			assert.Less(t, i, len(meta.Args)-1)
		case LenSourcePrev:
			assert.Greater(t, i, 0)
		case LenSourceReturnValue:
			assert.NotEqual(t, ArgTypeErrorCode, meta.ReturnValue.Type)
		}
	}
}
