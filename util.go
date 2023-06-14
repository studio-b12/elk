package elk

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

type errorJsonModel struct {
	Error   string    `json:"error"`
	Code    ErrorCode `json:"code,omitempty"`
	Message string    `json:"message,omitempty"`
	Details any       `json:"details,omitempty"`
}

// Json takes an error and marhals it into
// a JSON byte slice.
//
// If err is a wrapped error, the inner error
// will be represented in the "error" field.
// Otherwise, the result of Error() on err will
// be represented in the "error" field. This
// does only apply though if exposeError is
// passed as true. By default, "error" will
// contain no information about the actual
// error to prevent unintented information
// leakage.
//
// If the err implements HasCode, the code
// of the error will be represented in the
// "code" field of the JSON result.
//
// If the err implements HasMessage, the
// JSON object will contain it as "message"
// field, if present.
//
// When the JSON marshal fails, an error is
// returned.
func Json(err error, exposeError ...bool) ([]byte, error) {
	var model errorJsonModel

	if len(exposeError) > 0 && exposeError[0] {
		if inner := errors.Unwrap(err); inner != nil {
			model.Error = inner.Error()
		} else {
			model.Error = err.Error()
		}
	} else {
		model.Error = "internal error"
	}

	if mErr, ok := err.(HasMessage); ok {
		model.Message = mErr.Message()
	}

	if cErr, ok := err.(HasCode); ok {
		model.Code = cErr.Code()
	}

	data, jErr := json.MarshalIndent(model, "", "  ")
	if jErr != nil {
		return nil, jErr
	}

	return data, nil
}

// MustJson is an alias for Json but panics when
// the call to Json returns an error.
func MustJson(err error) []byte {
	return mustV(Json(err))
}

// JsonString behaves the same as Json() but returns the result as string instead
// of a slice of bytes.
func JsonString(err error, exposeError ...bool) (string, error) {
	res, err := Json(err, exposeError...)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// MustJson is an alias for JsonString but panics when the call to Json returns an error.
func MustJsonString(err error) string {
	return mustV(JsonString(err))
}

func mustV[TV any](v TV, err error) TV {
	if err != nil {
		panic(err)
	}
	return v
}
