package whoops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTypeOf(t *testing.T) {
	assert.True(t,
		IsOfType[testStringError](testStringError("test")))
	assert.True(t,
		IsOfType[testStructError](testStructError{"test"}))
	assert.True(t,
		IsOfType[testStringError](InnerError{Inner: testStringError("test")}))
	assert.True(t,
		IsOfType[*testRefError](InnerError{Inner: &testRefError{}}))

	assert.False(t,
		IsOfType[testStructError](testStringError("test")))
	assert.False(t,
		IsOfType[testStringError](testStructError{"test"}))
	assert.False(t,
		IsOfType[testStructError](InnerError{Inner: testStringError("test")}))
}

type testStringError string

func (t testStringError) Error() string { return string(t) }

type testStructError struct {
	msg string
}

func (t testStructError) Error() string { return t.msg }

type testRefError struct{}

func (t *testRefError) Error() string { return "refError" }
