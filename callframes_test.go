package elk

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/studio-b12/elk/internal/assert"
)

func TestFrames(t *testing.T) {
	stack := newCallStack(1, 3)
	assert.Equal(t, stack.offset, 1)
	assert.Equal(t, 3, len(stack.ptrs))
	assert.Nil(t, stack.frames)

	frames := stack.Frames()
	assert.Equal(t, 2, len(frames))
	assert.Equal(t, 3, len(stack.frames))
}

// This test checks what happens when the creation of a call
// stack happens in an anonymous function, which is then cleaned
// up by the garbage collector. Afterwards, Frames() is called on
// the CallStack to resolve the runtime.Frame objects from the
// internal slice of stack pointers.
func Test_stackCaptureGC(t *testing.T) {

	var cs *CallStack
	var cleanedUp bool

	{
		f := func() {
			cs = newCallStack(0, 10)
		}

		runtime.SetFinalizer(&f, func(a any) {
			cleanedUp = true
		})

		f()
	}

	// GC gets called twice to first call the registered
	// finalizer and then actually clean up the object.
	// See the documentation of runtime.SetFinalizer for
	// more information.
	runtime.GC()
	runtime.GC()

	assert.True(t, cleanedUp)
	for _, frame := range cs.Frames() {
		fmt.Printf("frame: %s\n", frame)
	}
}
