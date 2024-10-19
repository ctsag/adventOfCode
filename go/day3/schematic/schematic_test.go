package schematic

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"strings"
	"testing"
)

type MockFileReader struct {
	mock.Mock
}

func (mockReader *MockFileReader) Open(path string) (io.ReadCloser, error) {
	args := mockReader.Called(path)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestTotalCalculationShould(t *testing.T) {

	t.Run("calculate total of schematic values", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = "467..114..\n...*......\n..35..633.\n......#..."

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _ := CalculateTotal(fileName, mockReader)
		expected := 467 + 35 + 633

		assert.Equal(t, expected, actual, "Did not calculate the schematic values total correctly")
	})

	t.Run("fails when unable to read file", func(t *testing.T) {
		const fileName = "test_input.txt"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader("")), errors.New("file open error"))

		_, err := CalculateTotal(fileName, mockReader)
		expected := "file open error"

		assert.EqualError(
			t,
			err,
			expected,
			"Did not fail when unable to read file",
		)
	})

}

func TestSchematicExtractionShould(t *testing.T) {

	t.Run("extract schematic", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = "ABC\nDEF"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _ := extractSchematic(fileName, mockReader)

		var expected [][]byte
		expected = append(expected, []byte{'A', 'B', 'C'})
		expected = append(expected, []byte{'D', 'E', 'F'})

		assert.Equal(t, expected, actual, "Did not extract schematic correctly")
	})

	t.Run("extract values", func(t *testing.T) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '*', '.', '.', '.', '.', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '3', '5', '.', '.', '6', '3', '3', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '.', '.', '.', '.', '#', '.', '.'})

		actual := extractValues(schematic)
		expected := []schematicValue{
			{"467", 0, 0},
			{"114", 0, 5},
			{"35", 2, 2},
			{"633", 2, 6},
		}

		assert.Equal(t, expected, actual, "Did not extract values correctly")
	})

}

func TestAdjacencyDeterminationShould(t *testing.T) {

	tests := []struct {
		value    schematicValue
		expected bool
	}{
		{schematicValue{"467", 0, 0}, false},
		{schematicValue{"114", 0, 5}, true},
		{schematicValue{"32", 1, 4}, true},
		{schematicValue{"56", 1, 8}, true},
		{schematicValue{"501", 2, 2}, true},
		{schematicValue{"12", 2, 6}, true},
		{schematicValue{"6", 3, 9}, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("determine adjacency for value %s", test.value.num), func(t *testing.T) {
			var schematic [][]byte
			schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '#'})
			schematic = append(schematic, []byte{'.', '.', '.', '.', '*', '3', '2', '.', '5', '6'})
			schematic = append(schematic, []byte{'.', '.', '5', '0', '1', '.', '1', '2', '.', '.'})
			schematic = append(schematic, []byte{'.', '.', '.', '.', '.', '.', '/', '.', '.', '6'})

			actual := isAdjacentToSymbols(test.value, schematic)
			expected := test.expected

			assert.Equal(t, expected, actual, "Did not determine if adjacent to symbols")
		})
	}

}
