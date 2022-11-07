package tracer

import (
	"fmt"
	"testing"
)

func TestTypeConversions(t *testing.T) {
	for i := argStartInternal + 1; i < argEndInternal; i++ {
		t.Run(fmt.Sprintf("ArgType %d", i), func(t *testing.T) {
			if getHandler(i) == nil {
				t.Errorf("No handler for type %d", i)
			}
		})
	}
}
