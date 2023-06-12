package whoops

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
)

// CallFrame is a type alias for runtime.Frame with additional formatting
// functionailty used in downstream functions.
type CallFrame runtime.Frame

func (t CallFrame) String() string {
	return fmt.Sprintf("%s", t)
}

func (t CallFrame) Format(s fmt.State, verb rune) {
	width, hasWidth := s.Width()

	switch verb {
	case 's':
		format := "%s %s:%d"
		if hasWidth {
			format = "%-" + strconv.Itoa(width) + "s\t%s:%d"
		}
		fmt.Fprintf(s, format, t.Function, t.File, t.Line)
	case 'v':
		fmt.Fprintf(s, "%v", runtime.Frame(t))
	}
}

// CallStack contains the list of called
// runtime.Frames in the call chain with
// an offset from which frames are
// reported.
type CallStack struct {
	ptrs   []uintptr
	frames []CallFrame
	offset int
}

// Frames returns the offset slice of called
// runtime.Frame's in the recorded call stack.
func (t *CallStack) Frames() []CallFrame {
	if t.frames == nil {
		t.fetchCallFrames()
	}

	if len(t.frames) < t.offset {
		return nil
	}

	return t.frames[t.offset:]
}

// WriteIndent is an alias for write with the given
// indent string attached before each line of output.
func (t *CallStack) WriteIndent(w io.Writer, max int, indent string) {
	frames := t.Frames()

	if max > 0 && len(frames) > max {
		frames = frames[:max]
	}

	maxLenFName := 0
	for _, frame := range frames {
		if l := len(frame.Function); l > maxLenFName {
			maxLenFName = l
		}
	}
	for _, frame := range frames {
		fmt.Fprintf(w, "%s%"+strconv.Itoa(maxLenFName)+"s\n", indent, frame)
	}
}

// Write formats the call stack into a table of called
// function and the file plus line number and writes
// the result into the writer w.
//
// max defines the number of stack frames which are
// printed starting from the original caller.
func (t *CallStack) Write(w io.Writer, max int) {
	t.WriteIndent(w, max, "")
}

// String returns the formatted output of the callstack
// as string.
func (t *CallStack) String() string {
	var b bytes.Buffer
	t.Write(&b, 0)
	return b.String()
}

// At returns the formatted call frame at the given position n
// if existent.
func (t *CallStack) At(n int) (s string, ok bool) {
	frames := t.Frames()

	if n >= len(frames) || n < 0 {
		return "", false
	}

	frame := frames[n]
	return frame.String(), true
}

// First is shorthand for At(0) and returns the first frame in
// the CallStack, if available.
func (t *CallStack) First() (s string, ok bool) {
	return t.At(0)
}

func newCallStack(offset int, n int) *CallStack {
	callerPtrs := make([]uintptr, n)
	nPtrs := runtime.Callers(2, callerPtrs)

	return &CallStack{
		ptrs:   callerPtrs[:nPtrs],
		offset: offset,
	}
}

func (t *CallStack) fetchCallFrames() {
	frameCursor := runtime.CallersFrames(t.ptrs)

	callFrames := make([]CallFrame, 0, len(t.ptrs))
	for {
		frame, more := frameCursor.Next()
		callFrames = append(callFrames, CallFrame(frame))
		if !more {
			break
		}
	}

	t.frames = callFrames
}
