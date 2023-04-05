package whoops_test

import (
	"errors"
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

func TestWithDetailsOfType(t *testing.T) {
	t.Run("first", func(t *testing.T) {
		err := errors.New("err")
		err = whoops.Wrap(err, "first wrapped string")
		err = whoops.Wrap(err, fooImpl{})
		err = whoops.Wrap(err, "mid wrapped string")
		err = whoops.Wrap(err, 123.456)
		err = whoops.Wrap(err, "last wrapped string")

		_, ok := whoops.DetailsOfType[fooImpl](err)
		assert.True(t, ok)

		_, ok = whoops.DetailsOfType[foo](err)
		assert.True(t, ok)

		_, ok = whoops.DetailsOfType[int](err)
		assert.False(t, ok)

		_, ok = whoops.DetailsOfType[*fooImpl](err)
		assert.False(t, ok)

		v, ok := whoops.DetailsOfType[string](err)
		assert.True(t, ok)
		assert.Equal(t, "last wrapped string", v)
	})

	t.Run("last", func(t *testing.T) {
		err := errors.New("err")
		err = whoops.Wrap(err, "first wrapped string")
		err = whoops.Wrap(err, fooImpl{})
		err = whoops.Wrap(err, "mid wrapped string")
		err = whoops.Wrap(err, 123.456)
		err = whoops.Wrap(err, "last wrapped string")

		_, ok := whoops.DetailsOfType[fooImpl](err, true)
		assert.True(t, ok)

		_, ok = whoops.DetailsOfType[foo](err, true)
		assert.True(t, ok)

		_, ok = whoops.DetailsOfType[int](err, true)
		assert.False(t, ok)

		_, ok = whoops.DetailsOfType[*fooImpl](err, true)
		assert.False(t, ok)

		v, ok := whoops.DetailsOfType[string](err, true)
		assert.True(t, ok)
		assert.Equal(t, "first wrapped string", v)
	})
}
