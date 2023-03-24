package whoops_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/whoops"
)

func TestIsTypeOf(t *testing.T) {
	assert.True(t,
		whoops.IsOfType[testStringError](testStringError("test")))
	assert.True(t,
		whoops.IsOfType[testStructError](testStructError{"test"}))
	assert.True(t,
		whoops.IsOfType[testStringError](whoops.InnerError{Inner: testStringError("test")}))
	assert.True(t,
		whoops.IsOfType[*testRefError](whoops.InnerError{Inner: &testRefError{}}))

	assert.False(t,
		whoops.IsOfType[testStructError](testStringError("test")))
	assert.False(t,
		whoops.IsOfType[testStringError](testStructError{"test"}))
	assert.False(t,
		whoops.IsOfType[testStructError](whoops.InnerError{Inner: testStringError("test")}))
}

type testStringError string

func (t testStringError) Error() string { return string(t) }

type testStructError struct {
	msg string
}

func (t testStructError) Error() string { return t.msg }

type testRefError struct{}

func (t *testRefError) Error() string { return "refError" }
