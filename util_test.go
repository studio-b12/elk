package elk_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/studio-b12/elk"
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
		elk.IsOfType[stringError](stringError("test")))
	assert.True(t,
		elk.IsOfType[structError](structError{"test"}))
	assert.True(t,
		elk.IsOfType[stringError](elk.InnerError{Inner: stringError("test")}))
	assert.True(t,
		elk.IsOfType[*refError](elk.InnerError{Inner: &refError{}}))

	assert.False(t,
		elk.IsOfType[structError](stringError("test")))
	assert.False(t,
		elk.IsOfType[stringError](structError{"test"}))
	assert.False(t,
		elk.IsOfType[structError](elk.InnerError{Inner: stringError("test")}))
}

type foo interface {
	Foo()
}

type fooImpl struct{}

func (fooImpl) Foo() {}
