package whoops

import (
	"encoding/json"
	"errors"
)

// UnwrapFull takes an error and unwraps it until
// it can not be unwrapped anymore. Then, the
// last error is returned.
func UnwrapFull(err error) error {
	for {
		uErr := errors.Unwrap(err)
		if uErr == nil {
			return err
		}
		err = uErr
	}
}

// As applies errors.As() on the given err
// using the given type T as target for the
// unwrapping.
//
// Refer to the documentation of errors.As()
// for more details:
// https://pkg.go.dev/errors#As
func As[T error](err error) (t T, ok bool) {
	ok = errors.As(err, &t)
	return t, ok
}

// IsOfType returns true when the given
// error is of the type of T.
//
// If not and the error can be unwrapped,
// the unwrapped error will be checked
// until it either matches the type T or
// can not be further unwrapped.
func IsOfType[T error](err error) bool {
	_, ok := err.(T)
	if ok {
		return true
	}

	err = errors.Unwrap(err)
	if err != nil {
		return IsOfType[T](err)
	}

	return false
}

// Format returns the formatted error message
// when err implements Formatted. Otherwise,
// the result of err.Error() is returned.
func Format(err error) string {
	if fErr, ok := err.(HasFormat); ok {
		return fErr.Formatted()
	}
	return err.Error()
}

// Message returns the message when error
// implements HasMessage. Otherwise, the
// content of err.Error() is returned.
func Message(err error) string {
	if mErr, ok := err.(HasMessage); ok {
		return mErr.Message()
	}
	return err.Error()
}

func Code(err error) ErrorCode {
	if cErr, ok := err.(HasCode); ok {
		return cErr.Code()
	}
	return CodeUnexpected
}

type errorJsonModel struct {
	Error   string    `json:"error,omitempty"`
	Code    ErrorCode `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
	Details any       `json:"details,omitempty"`
}

// Json takes an error and marhals it into
// a JSON string. It will contain the
// inner error as "error" field.
//
// If the err implements HasMessage, the
// JSON object will contain it as "message"
// field, if present.
//
// If the err implements HasDetails, the
// JSON object will contain them as "details"
// field, if present. If details is an error,
// the contents of the Error() method are
// represented in the field.
//
// When the JSON marshal fails, an error is
// returned.
func Json(err error, includeInner ...bool) (string, error) {
	var model errorJsonModel

	if len(includeInner) > 0 && includeInner[0] {
		if inner := errors.Unwrap(err); inner != nil {
			model.Error = inner.Error()
		} else {
			model.Error = err.Error()
		}
	}

	if mErr, ok := err.(HasMessage); ok {
		model.Message = mErr.Message()
	}

	if cErr, ok := err.(HasCode); ok {
		model.Code = cErr.Code()
	}

	data, jErr := json.MarshalIndent(model, "", "  ")
	if jErr != nil {
		return "", jErr
	}

	return string(data), nil
}

// MustJson is an alias for Json but panics when
// the call to Json returns an error.
func MustJson(err error) string {
	return mustV(Json(err))
}

func mustV[TV any](v TV, err error) TV {
	if err != nil {
		panic(err)
	}
	return v
}
