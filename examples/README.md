# examples

Here you can find a conclusive example on how elk can be used to control error behavior in a very
simple REST API application with a database, a controller and corresponding route handlers.

VBelow you can see a visualization on how an unexpected database error propagates through the different layers of the application.

![](../.github/media/example-flow.png)

## Setup

You can simply run the example by executing the following command from the root of the repository.
```
go run examples/server/*.go
```

After that, you can call the endpoints `GET http://localhost:8080/count?id=1` and 
`POST http://localhost:8080/count?id=1`. When no "database" (the `db.json` file) is initialized, the
first calls to these endpoints will intentionally fail with a `500 Internal Server Error` and the details
of that error are logged to the console. Therefore, the `GET` endpoint demonstrates the detailed error
representation (formatted with `%+v`) and the `POST` endpoint shows the verbose error representation
(formatted with `%#v`).

After the first call to one of those endpoints, the `db.json` is created in the execution directory
and subsequent requests will succeed.

If an ID is passed to the `GET` endpoint which is not in the "database", the endpoint will return a
`404 Not Found` error though.