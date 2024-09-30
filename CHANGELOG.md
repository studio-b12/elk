# Changelog

## v0.4.0

- Updated [`Cast`](https://pkg.go.dev/github.com/studio-b12/elk#Cast) to make use of `errors.As` to unwrap errors.
- Added [`ErrorResponseModel`](https://pkg.go.dev/github.com/studio-b12/elk#ErrorResponseModel) and [`ToErrorResponseModel`](https://pkg.go.dev/github.com/studio-b12/elk#ToErrorResponseModel).

## v0.3.0

- `Cast` does now handle joined errors created via [`errors.Join`](https://pkg.go.dev/errors#Join). See [the documentation](https://pkg.go.dev/github.com/studio-b12/elk@v0.3.0#Cast) for more information.
- Fixed some typos in the code documentation.

## v0.2.0

- Added function `WrapCopyCode`, which takes an error and–if it has an error code–copies the code to the new wrapped error.
- Added formatted function alternatives `NewErrorf`, `Wrapf` and `WrapCopyCodef`

## v0.1.0

- Initial preview release.
