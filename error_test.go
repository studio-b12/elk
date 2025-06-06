package elk

import (
	"errors"
	"testing"

	"github.com/studio-b12/elk/internal/assert"
)

type customModelWithCode struct {
	code ErrorCode
}

func (t *customModelWithCode) Error() string   { return "customModelWithCode" }
func (t *customModelWithCode) Code() ErrorCode { return t.code }

type customModelWithCodeAndMessage struct {
	code    ErrorCode
	message string
}

func (t *customModelWithCodeAndMessage) Error() string   { return "customModelWithCodeAndMessage" }
func (t *customModelWithCodeAndMessage) Code() ErrorCode { return t.code }
func (t *customModelWithCodeAndMessage) Message() string { return t.message }

func TestCast(t *testing.T) {
	const ErrCode = ErrorCode("some-error-code")

	t.Run("error", func(t *testing.T) {
		err := errors.New("some error")
		wrappedErr := Wrap(ErrCode, err)
		castError := Cast(wrappedErr)

		assert.Equal(t, wrappedErr, castError)
		assert.Equal(t, err, castError.Unwrap())
	})

	t.Run("error-default", func(t *testing.T) {
		err := errors.New("some error")
		castError := Cast(err)

		assert.Equal(t, castError.Code(), CodeUnexpected)
		assert.Equal(t, err, castError.Unwrap())
	})

	t.Run("error-custom", func(t *testing.T) {
		err := errors.New("some error")
		castError := Cast(err, ErrCode)

		assert.Equal(t, castError.Code(), ErrCode)
		assert.Equal(t, err, castError.Unwrap())
	})

	t.Run("error-join", func(t *testing.T) {
		err := errors.Join(nil)
		castError := Cast(err, ErrCode)
		assert.Equal(t, castError.Code(), ErrCode)

		err = errors.Join(errors.New("foo"), errors.New("bar"))
		castError = Cast(err, ErrCode)
		assert.Equal(t, castError.Code(), ErrCode)

		customCode := ErrorCode("custom-code")
		customCode2 := ErrorCode("custom-code-2")

		err = errors.Join(errors.New("foo"), NewError(customCode))
		castError = Cast(err, ErrCode)
		assert.Equal(t, castError.Code(), customCode)

		err = errors.Join(NewError(customCode), NewError(customCode2))
		castError = Cast(err, ErrCode)
		assert.Equal(t, castError.Code(), ErrCode)
	})

	t.Run("custom-model", func(t *testing.T) {
		errCode := ErrorCode("custom-code")

		type MyError struct {
			InnerError
			SomeData string
		}

		err := MyError{
			InnerError: InnerError{Inner: NewError(errCode, "some message")},
			SomeData:   "some data",
		}

		castErrCode := Cast(err).Code()
		assert.Equal(t, CodeUnexpected, castErrCode)
	})

	t.Run("custom-model-wrapped", func(t *testing.T) {
		errCode := ErrorCode("custom-code")
		wrappedErrCode := ErrorCode("wrapped-custom-code")

		type MyError struct {
			InnerError
			SomeData string
		}

		wrappedErr := NewError(wrappedErrCode, "some message")

		err := MyError{
			InnerError: InnerError{Inner: Wrap(errCode, wrappedErr)},
			SomeData:   "some data",
		}

		castErrCode := Cast(err).Code()
		assert.Equal(t, CodeUnexpected, castErrCode)
	})

	t.Run("custom-model-with-code", func(t *testing.T) {
		errCode := ErrorCode("custom-code")
		err := &customModelWithCode{code: errCode}

		castErrCode := Cast(err).Code()
		assert.Equal(t, errCode, castErrCode)
	})

	t.Run("custom-model-with-code-and-message", func(t *testing.T) {
		errCode := ErrorCode("custom-code")
		errMessage := "custom error message"
		err := &customModelWithCodeAndMessage{code: errCode, message: errMessage}

		castErr := Cast(err)

		assert.Equal(t, errCode, castErr.Code())
		assert.Equal(t, errMessage, castErr.Message())
	})
}
