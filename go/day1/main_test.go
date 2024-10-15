package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestOutputsTotal(t *testing.T) {
	const filename = "testdata/test_input.txt"
	const expectedTotal = 19 + 66 + 78 + 29 + 83 + 14 + 76

	os.Args = []string{"cmd", filename}

	actualOut, actualCode, err := captureStdOut(main)
	if err != nil {
		assert.Fail(t, "Error: %s\n", err)
	}

	expectedOut := fmt.Sprintf("The sum of all calibration values is %d\n", expectedTotal)
	expectedCode := 0

	assert.Equal(
		t,
		expectedCode,
		actualCode,
		"Did not exit with the expected code",
	)

	assert.Equal(
		t,
		expectedOut,
		actualOut,
		"Did not output the total correctly",
	)
}

func TestFailsWhenWrongArgs(t *testing.T) {
	os.Args = []string{"cmd"}

	actualCode := captureErrorCode(main)
	expectedCode := 1

	assert.Equal(
		t,
		expectedCode,
		actualCode,
		"Did not exit with the expected code",
	)
}

func TestFailsForFileErrors(t *testing.T) {
	const filename = "non_existent_file.txt"

	os.Args = []string{"cmd", filename}

	actualCode := captureErrorCode(main)
	expectedCode := 2

	assert.Equal(
		t,
		expectedCode,
		actualCode,
		"Did not exit with the expected code",
	)
}

func captureStdOut(f func()) (string, int, error) {
	originalStdout := os.Stdout
	originalExit := osExit
	exitCode := 0

	defer func() {
		os.Stdout = originalStdout
		osExit = originalExit
	}()

	osExit = func(code int) {
		exitCode = code
	}

	inPipe, outPipe, _ := os.Pipe()
	os.Stdout = outPipe

	f()

	if outPipe.Close() != nil {
		return "", -1, errors.New("unable to close output pipe")
	}

	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, inPipe)
	if err != nil {
		return "", -1, errors.New("unable to capture input pipe")
	}

	return buffer.String(), exitCode, nil
}

func captureErrorCode(f func()) int {
	originalExit := osExit
	exitCode := 0

	defer func() {
		osExit = originalExit
	}()

	osExit = func(code int) {
		exitCode = code
	}

	f()

	return exitCode
}
