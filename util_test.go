package whoops_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/whoops"
)

type stringError string

func (t stringError) Error() string { return string(t) }

type structError struct {
	msg string
}

func (t structError) Error() string { return t.msg }

type refError struct{}

func (t *refError) Error() string { return "refError" }

func TestIsTypeOf(t *testing.T) {
	assert.True(t,
		whoops.IsOfType[stringError](stringError("test")))
	assert.True(t,
		whoops.IsOfType[structError](structError{"test"}))
	assert.True(t,
		whoops.IsOfType[stringError](whoops.InnerError{Inner: stringError("test")}))
	assert.True(t,
		whoops.IsOfType[*refError](whoops.InnerError{Inner: &refError{}}))

	assert.False(t,
		whoops.IsOfType[structError](stringError("test")))
	assert.False(t,
		whoops.IsOfType[stringError](structError{"test"}))
	assert.False(t,
		whoops.IsOfType[structError](whoops.InnerError{Inner: stringError("test")}))
}

type foo interface {
	Foo()
}

type fooImpl struct{}

func (fooImpl) Foo() {}
