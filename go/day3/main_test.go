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

func TestApplicationShould(t *testing.T) {

	t.Run("output the schematic values total", func(t *testing.T) {
		const filename = "testdata/test_input.txt"
		const expectedTotal = 467 + 35 + 633 + 617 + 592 + 755 + 664 + 598

		os.Args = []string{"cmd", filename}

		actualOut, actualCode, _ := captureStdOut(main)

		expectedOut := fmt.Sprintf("The sum of all schematic values is %d\n", expectedTotal)
		expectedCode := 0

		assert.Equal(t, expectedCode, actualCode, "Did not exit with the expected code")
		assert.Contains(t, actualOut, expectedOut, "Did not output the schematic values total correctly")
	})

	t.Run("output the gear ratios total", func(t *testing.T) {
		const filename = "testdata/test_input.txt"
		const expectedTotal = 467*35 + 755*598

		os.Args = []string{"cmd", filename}

		actualOut, actualCode, _ := captureStdOut(main)

		expectedOut := fmt.Sprintf("The sum of all gear ratios is %d\n", expectedTotal)
		expectedCode := 0

		assert.Equal(t, expectedCode, actualCode, "Did not exit with the expected code")
		assert.Contains(t, actualOut, expectedOut, "Did not output the gear ratios total correctly")
	})

	t.Run("fail when wrong arguments are passed", func(t *testing.T) {
		os.Args = []string{"cmd"}

		actualCode := captureErrorCode(main)
		expectedCode := 1

		assert.Equal(t, expectedCode, actualCode, "Did not exit with the expected code")
	})

	t.Run("fail for any file parsing error", func(t *testing.T) {
		os.Args = []string{"cmd", "non_existent_file.txt"}

		actualCode := captureErrorCode(main)
		expectedCode := 2

		assert.Equal(t, expectedCode, actualCode, "Did not exit with the expected code")
	})

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
