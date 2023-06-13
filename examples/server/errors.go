package main

import "github.com/studio-b12/whoops"

const (
	ErrorInternal      = whoops.ErrorCode("internal-server-error")
	ErrorCountNotFound = whoops.ErrorCode("cound-not-found")
)
