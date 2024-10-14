package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidationShould(t *testing.T) {

	t.Run("extract single argument", func(t *testing.T) {
		var args = []string{"cmd", "input.txt"}

		actual, _ := ExtractSingleArgIgnoringOthers(args, 2)
		expected := "input.txt"

		assert.Equal(
			t,
			expected,
			actual,
			"Did not extract path argument",
		)
	})

	t.Run("ignore extraneous arguments", func(t *testing.T) {
		var args = []string{"cmd", "input.txt", "extraneous"}

		actual, _ := ExtractSingleArgIgnoringOthers(args, 2)
		expected := "input.txt"

		assert.Equal(
			t,
			expected,
			actual,
			"Did not extract path argument",
		)
	})

	t.Run("fail for insufficient arguments", func(t *testing.T) {
		var args = []string{"cmd"}

		_, err := ExtractSingleArgIgnoringOthers(args, 2)
		expected := "no file parameter provided"

		assert.Error(
			t,
			err,
			expected,
		)
	})

}
