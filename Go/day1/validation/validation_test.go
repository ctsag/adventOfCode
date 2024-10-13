package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestExtractsSingleArg(t *testing.T) {
	var args = []string{"cmd", "input.txt"}

	actual, _ := ExtractSingleArgIgnoringOthers(args, 2)
	expected := "input.txt"

	assert.Equal(
		t,
		expected,
		actual,
		"Did not extract path argument",
	)
}

func TestIgnoresExtraneousArguments(t *testing.T) {
	var args = []string{"cmd", "input.txt", "extraneous"}

	actual, _ := ExtractSingleArgIgnoringOthers(args, 2)
	expected := "input.txt"

	assert.Equal(
		t,
		expected,
		actual,
		"Did not extract path argument",
	)
}

func TestFailsWhenNotEnoughArguments(t *testing.T) {
	var args = []string{"cmd"}

	_, err := ExtractSingleArgIgnoringOthers(args, 2)
	expected := "no file parameter provided"

	assert.Error(
		t,
		err,
		expected,
	)
}
