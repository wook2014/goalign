package io

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

const (
	EXIT_SUCCESS = 0
	EXIT_FAILURE = 1
)

func ExitWithMessage(err error) {
	_, fn, line, _ := runtime.Caller(1)

	name := strings.Split(fn, "/goalign/")[1]
	fmt.Fprintf(os.Stderr, "[Error] in %s (line %d), message: %v\n", name, line, err)
	os.Exit(EXIT_FAILURE)
}

func PrintMessage(message string) {
	_, fn, line, _ := runtime.Caller(1)

	name := strings.Split(fn, "/goalign/")[1]
	fmt.Fprintf(os.Stderr, "[Warning] in %s (line %d), message: %s\n", name, line, message)
}

func PrintSimpleMessage(message string) {
	fmt.Fprintf(os.Stderr, "%s\n", message)
}

func LogError(err error) {
	_, fn, line, _ := runtime.Caller(1)
	name := strings.Split(fn, "/goalign/")[1]
	fmt.Fprintf(os.Stderr, "[Error] in %s (line %d), message: %v\n", name, line, err)
}
