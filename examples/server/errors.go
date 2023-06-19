package main

import "github.com/studio-b12/elk"

const (
	ErrorInternal      = elk.ErrorCode("internal-server-error")
	ErrorCountNotFound = elk.ErrorCode("count-not-found")
)
