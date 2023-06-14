package main

import (
	"errors"
	"fmt"

	"github.com/studio-b12/elk"
)

func main() {
	const MyErrorCode = elk.ErrorCode("my-error-code")

	{
		fmt.Println("Simple formatting example (with message):")

		err := elk.Wrap(MyErrorCode,
			errors.New("somethign went wrong"),
			"Damn, what happened?")
		fmt.Printf("%s\n", err)
	}

	{
		fmt.Println("\nSimple formatting example (without message):")

		err := elk.Wrap(MyErrorCode,
			errors.New("somethign went wrong"))
		fmt.Printf("%s\n", err)
	}

	{
		fmt.Println("\nSimple detailed formatting example (with message):")

		err := elk.Wrap(MyErrorCode,
			errors.New("somethign went wrong"),
			"Damn, what happened?")
		fmt.Printf("%v\n", err)
	}

	{
		fmt.Println("\nSimple detailed formatting example (without message):")

		err := elk.Wrap(MyErrorCode,
			errors.New("somethign went wrong"))
		fmt.Printf("%v\n", err)
	}

	{
		fmt.Println("\nMore detailed formatting example:")

		err := elk.Wrap(MyErrorCode,
			errors.New("somethign went wrong"),
			"Damn, what happened?")
		fmt.Printf("%+.5v\n", err)
	}

	{
		fmt.Println("\nVerbose formatting example:")

		err := elk.Wrap(MyErrorCode,
			errors.New("somethign went wrong"),
			"Damn, what happened?")
		fmt.Printf("%#v\n", err)
	}
}
