package elk

import (
	"errors"
	"testing"

	"github.com/studio-b12/elk/internal/assert"
)

func TestCast(t *testing.T) {
	const ErrCode = ErrorCode("some-error-code")

	t.Run("cast-Error", func(t *testing.T) {
		err := errors.New("some error")
		wrappedErr := Wrap(ErrCode, err)
		castError := Cast(wrappedErr)

		assert.Equal(t, wrappedErr, castError)
		assert.Equal(t, err, castError.Unwrap())
	})

	t.Run("cast-error-default", func(t *testing.T) {
		err := errors.New("some error")
		castError := Cast(err)

		assert.Equal(t, castError.Code(), CodeUnexpected)
		assert.Equal(t, err, castError.Unwrap())
	})

	t.Run("cast-error-custom", func(t *testing.T) {
		err := errors.New("some error")
		castError := Cast(err, ErrCode)

		assert.Equal(t, castError.Code(), ErrCode)
		assert.Equal(t, err, castError.Unwrap())
	})
}
