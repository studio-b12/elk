# 🦌 elk

An extensive error package with focus on comprehensiveness, tracability and ease of use.

## Getting started

elk provides a simple error model called `Error`. It is classified by an `ErrorCode` and either wraps a given inner error or creates one if there is no underlying error. You can also pass an optional message to provide more detailed context to the error. `Error`s also record the callstack from where they have been created so that they can be easily traced thorugh the codebase, if necessary.

Create a detailed error witn an error code and message.
```go
const ErrDeviceNotFound = elk.ErrorCode("device-not-found")

err := elk.Detailed(ErrDeviceNotFound, "the device could not be found")
```

Wrap a previous error with an error code and message.
```go
device, err := db.GetDevice(id)
if err != nil {
    err = elk.Detailed(elk.CodeUnexpected, 
        "failed receiving device from database")
}
```

`Error` also implements the [`fmt.Formatter`](https://pkg.go.dev/fmt#Formatter) interface so you can finely control how errors are displayed. See the [Formatting](#formatting) section for more information.

The recommended way to use this construct is to wrap an error on each layer in your application where the error changes the state of the outcome of the error. In example, when your database returns an `ErrNoRows` error and in your controller, that means that no values could be found for the given request, you can wrap the original database error with an error Code (`ErrObjectNotFound` i.E.) and an additional message to clarify what went wrong to either the user or developers of the layers above, if desired.

This way, you can give other meaning to errors on each layer without losign details about each consecutive error.

### Formatting

## Contribute

If you find any issues, want to submit a suggestion for a new feature or improvement of an existing one or just want to ask a question, feel free to [create an Issue](https://github.com/studio-b12/elk/issues/new).

If you want to contribute to the project, just [create a fork](https://github.com/studio-b12/elk/fork) and [create a pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request) with your changes. We are happy to review your contribution and make you a part of the project. 😄

---

© 2023 B12-Touch GmbH  
https://b12-touch.de

Covered by the MIT License.