# 🦌 elk

An extensive error package with focus on comprehensiveness, tracability and ease of use.

## Getting started

elk provides a simple error model called `Error`. It is classified by an `ErrorCode` and either wraps a given inner error or creates one if there is no underlying error. You can also pass an optional message to provide more detailed context to the error. `Error`s also record the callstack from where they have been created so that they can be easily traced thorugh the codebase, if necessary.

Create a detailed error witn an error code and message.
```go
const ErrDeviceNotFound = elk.ErrorCode("device-not-found")

err := elk.NewError(ErrDeviceNotFound, "the device could not be found")
```

Wrap a previous error with an error code and message.
```go
device, err := db.GetDevice(id)
if err != nil {
    err = elk.Wrap(elk.CodeUnexpected, err,
        "failed receiving device from database")
}
```

`Error` also implements the [`fmt.Formatter`](https://pkg.go.dev/fmt#Formatter) interface so you can finely control how errors are displayed. See the [Formatting](#formatting) section for more information.

The recommended way to use this construct is to wrap an error on each layer in your application where the error changes the state of the outcome of the error. In example, when your database returns an `ErrNoRows` error and in your controller, that means that no values could be found for the given request, you can wrap the original database error with an error Code (`ErrObjectNotFound` i.E.) and an additional message to clarify what went wrong to either the user or developers of the layers above, if desired.

This way, you can give other meaning to errors on each layer without losign details about each consecutive error.

### Formatting

> In [`examples/formatting`](examples/formatting), you can find the different formatting options in use. Execute it to see them in action in your terminal!

As mentioned above, `Error` implements [`fmt.Formatter`](https://pkg.go.dev/fmt#Formatter). So there are some custom options for printing `Error` instances.

#### `%s` or `%q`

When printing the error as a simple string, it will reproduce an output simmilar to Go's default error types with a single message in a single line. If the error has a message, the message is snown. Otherwise, the `%s` formatted contents of the inner error is displayed.

```go
const MyErrorCode = elk.ErrorCode("my-error-code")

err := elk.Wrap(MyErrorCode,
	errors.New("somethign went wrong"),
	"Damn, what happened?")

fmt.Printf("%s\n", err)
// Output: Damn, what happened?
```

#### `%v`

Without any further flags, this prints a single line combined output of the wrapped errors code, message (if set) and inner errors text.

```go
const MyErrorCode = elk.ErrorCode("my-error-code")

err := elk.Wrap(MyErrorCode,
	errors.New("somethign went wrong"),
	"Damn, what happened?")

fmt.Printf("%v\n", err)
// Output: <my-error-code> Damn, what happened? (somethign went wrong)
```

With the additional flag `+`, more details are shown like the callstack (see [Callstack secion](#callstack)) of the error and the inner error. By passing the precision parameter (i.E. `%+.5v`), you can specify the maximum depth of the shown callstack. By default, a depth of `1000` is assumed. If you set this to `0`, no call stack is printed.

```go
const MyErrorCode = elk.ErrorCode("my-error-code")

err := elk.Wrap(MyErrorCode,
	errors.New("somethign went wrong"),
	"Damn, what happened?")

fmt.Printf("%+.5v\n", err)
// Output:
// <my-error-code> Damn, what happened?
// stack:
//   main.main             /home/r.hoffmann@intern.b12-group.de/dev/lib/whoops/examples/formatting/main.go:50
//   runtime.main          /home/r.hoffmann@intern.b12-group.de/.local/goup/current/go/src/runtime/proc.go:250
//   runtime.goexit        /home/r.hoffmann@intern.b12-group.de/.local/goup/current/go/src/runtime/asm_amd64.s:1598
// inner error:
//   somethign went wrong
```

By setting the flag `#`, you can enable a verbose view of the error. This unwraps all layers of the error and prints a detailed overview of each visted error containing the error string, origin (where it has been wrapped) and the type of the error. You can also specify the maximum depth that shall be displayed by giving the precision parameter (i.E. `%#.5v`). When not specified, a default value of `1000` is assumed.

```go
const MyErrorCode = elk.ErrorCode("my-error-code")

err := elk.Wrap(MyErrorCode,
	errors.New("somethign went wrong"),
	"Damn, what happened?")

fmt.Printf("%#.5v\n", err)
// Output:
// <my-error-code> Damn, what happened?
// originated:
//   main.main /home/r.hoffmann@intern.b12-group.de/dev/lib/whoops/examples/formatting/main.go:59
// type:
//   elk.Error
// ----------
// somethign went wrong
// type:
//   *errors.errorString
// ----------
```

### Callstack

When creating an `Error`–either by wrapping a previous error using `Wrap` or creating it using `NewError`–, it records where it has been wrapped in the Code in a `CallStack` object. This can then be accessed via the `CallStack` getter or is displayed when using the detailed and verbose formatting options as shown previously.

The `CallStack` contains a list of subsequent callers starting from the point where the `CallStack` has been created (when creating an `Error` instance, i.E.) followed by each previous caller of that function.

This `CallStack` object efficiently stores the frame pointers and resolves the context when calling the `Frames` getter on it.

Inner frames are wrapped using the `CallFrame` type, which also provides some formatting utilities.

Using the `%s` formatting verb, the `CallFrame` is printed in the following format.
```
main.main /home/me/dev/lib/elk/examples/formatting/main.go:59
```

When using the `%v` verb, it is formatted using the `%v` formatting on the underlying `runtime.Frame`.

## Contribute

If you find any issues, want to submit a suggestion for a new feature or improvement of an existing one or just want to ask a question, feel free to [create an Issue](https://github.com/studio-b12/elk/issues/new).

If you want to contribute to the project, just [create a fork](https://github.com/studio-b12/elk/fork) and [create a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) with your changes. We are happy to review your contribution and make you a part of the project. 😄

---

© 2023 B12-Touch GmbH  
https://b12-touch.de

Covered by the MIT License.