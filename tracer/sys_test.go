package tracer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SyscallSupport(t *testing.T) {
	for i := 0; i <= 335; i++ {
		t.Run(fmt.Sprintf("syscall %d", i), func(t *testing.T) {
			meta, err := LookupSyscall(i)
			require.NoError(t, err)
			require.NotNil(t, meta)
			assert.NotEmpty(t, meta.Name)
		})
	}
}
