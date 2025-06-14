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

// IsCode is shorthand for `elk.Cast(err).Code() == errorCode`.
func IsCode(err error, code ErrorCode) bool {
	return Cast(err).Code() == code
}

// ErrorResponseModel is used to encode an Error into an API response.
type ErrorResponseModel struct {
	Code    ErrorCode // The error code
	Message string    `json:",omitempty"` // An optional short message to further specify the error
	Status  int       `json:",omitempty"` // An optional platform- or protocol-specific status code; i.e. HTTP status code
	Details any       `json:",omitempty"` // Optional additional detailed context for the error
}

// ToResponseModel transforms the
func (t Error) ToResponseModel(statusCode int) (model ErrorResponseModel) {
	model.Status = statusCode
	model.Code = t.Code()

	if mErr, ok := As[HasMessage](t); ok {
		model.Message = mErr.Message()
	}

	if dErr, ok := As[HasDetails](t); ok {
		model.Details = dErr.Details()
	}

	return model
}

// Json takes an error and marshals it into
// a JSON byte slice.
//
// If err is a wrapped error, the inner error
// will be represented in the "error" field.
// Otherwise, the result of Error() on err will
// be represented in the "error" field. This
// does only apply though if exposeError is
// passed as true. By default, "error" will
// contain no information about the actual
// error to prevent unintended information
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
func Json(err error, statusCode int) ([]byte, error) {
	model := Cast(err).ToResponseModel(statusCode)

	data, jErr := json.MarshalIndent(model, "", "  ")
	if jErr != nil {
		return nil, jErr
	}

	return data, nil
}

// MustJson is an alias for Json but panics when
// the call to Json returns an error.
func MustJson(err error, statusCode int) []byte {
	return mustV(Json(err, statusCode))
}

// JsonString behaves the same as Json() but returns the result as string instead
// of a slice of bytes.
func JsonString(err error, statusCode int) (string, error) {
	res, err := Json(err, statusCode)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// MustJsonString is an alias for JsonString but panics when the call to Json returns an error.
func MustJsonString(err error, statusCode int) string {
	return mustV(JsonString(err, statusCode))
}

func mustV[TV any](v TV, err error) TV {
	if err != nil {
		panic(err)
	}
	return v
}
