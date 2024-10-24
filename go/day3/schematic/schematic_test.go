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

		actual, _, _ := CalculateTotals(fileName, mockReader)
		expected := 467 + 35 + 633

		assert.Equal(t, expected, actual, "Did not calculate the schematic values total correctly")
	})

	t.Run("calculate total of gear ratios", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = "467..114..\n...*......\n..35..633.\n......#...\n617*......\n..58......"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		_, actual, _ := CalculateTotals(fileName, mockReader)
		expected := 467*35 + 617*58

		assert.Equal(t, expected, actual, "Did not calculate the gear ratios total correctly")
	})

	t.Run("fails when unable to read file", func(t *testing.T) {
		const fileName = "test_input.txt"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader("")), errors.New("file open error"))

		_, _, err := CalculateTotals(fileName, mockReader)
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

	t.Run("extract schematic values", func(t *testing.T) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '*', '.', '.', '.', '.', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '3', '5', '.', '.', '6', '3', '3', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '.', '.', '.', '.', '#', '.', '.'})

		actual, _ := extractTokens(schematic)
		expected := []schematicValue{
			{"467", 0, 0},
			{"114", 0, 5},
			{"35", 2, 2},
			{"633", 2, 6},
		}

		assert.Equal(t, expected, actual, "Did not extract schematic values correctly")
	})

	t.Run("extract gears", func(t *testing.T) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '*', '.', '.', '.', '*', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '3', '5', '.', '.', '.', '3', '3', '.'})
		schematic = append(schematic, []byte{'.', '*', '.', '.', '.', '*', '.', '#', '.', '.'})

		_, actual := extractTokens(schematic)
		expected := []gear{
			{1, 3},
			{1, 7},
			{3, 1},
			{3, 5},
		}

		assert.Equal(t, expected, actual, "Did not extract gears correctly")
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

func TestRatioDeterminationShould(t *testing.T) {

	tests := []struct {
		value    gear
		expected int
	}{
		{gear{0, 9}, 0},
		{gear{1, 1}, 467 * 501},
		{gear{1, 4}, 0},
		{gear{2, 9}, 56 * 6},
		{gear{3, 0}, 0},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("determine ratio for gear at %v", test.value), func(t *testing.T) {
			var schematic [][]byte
			schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '*'})
			schematic = append(schematic, []byte{'.', '*', '.', '.', '*', '3', '2', '.', '5', '6'})
			schematic = append(schematic, []byte{'.', '.', '5', '0', '1', '.', '1', '2', '.', '*'})
			schematic = append(schematic, []byte{'*', '.', '.', '.', '.', '.', '/', '.', '.', '6'})

			actual := getGearRatioOrZero(test.value, schematic)
			expected := test.expected

			assert.Equal(t, expected, actual, "Did not determine gear ratio correctly")
		})
	}

}

func TestValueFromCoordinates(t *testing.T) {

	tests := []struct {
		row      int
		col      int
		expected string
		isError  bool
	}{
		{0, 0, "467", false},
		{1, 5, "32", false},
		{1, 9, "56", false},
		{2, 3, "501", false},
		{2, 6, "1", false},
		{3, 9, "6", false},
		{0, 3, "", true},
		{0, 9, "", true},
		{2, 5, "", true},
		{3, 0, "", true},
		{3, 8, "", true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("get value from coordinates %d,%d", test.row, test.col), func(t *testing.T) {
			var schematic [][]byte
			schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '.', '.', '.', '.', '#'})
			schematic = append(schematic, []byte{'.', '.', '.', '.', '*', '3', '2', '.', '5', '6'})
			schematic = append(schematic, []byte{'.', '.', '5', '0', '1', '.', '1', '.', '.', '.'})
			schematic = append(schematic, []byte{'.', '.', '.', '.', '.', '.', '/', '.', '.', '6'})

			actual, err := getValueFromCoordinates(test.row, test.col, schematic)
			expected := test.expected
			isError := err != nil

			assert.Equal(t, test.isError, isError, "Did not get value from coordinates")
			assert.Equal(t, expected, actual.num, "Did not get correct value from coordinates")
		})
	}

}

func BenchmarkTotalsCalculation(b *testing.B) {

	b.Run("totals calculation", func(b *testing.B) {
		const fileName = "test_input.txt"
		const lines = "467..114..\n...*......\n..35..633.\n......#..."

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		for i := 0; i < b.N; i++ {
			_, _, _ = CalculateTotals(fileName, mockReader)
		}
	})

	b.Run("schematic extraction", func(b *testing.B) {
		const fileName = "test_input.txt"
		const lines = "ABC\nDEF"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		for i := 0; i < b.N; i++ {
			_, _ = extractSchematic(fileName, mockReader)
		}
	})

	b.Run("schematic values extraction", func(b *testing.B) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '*', '.', '.', '.', '.', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '3', '5', '.', '.', '6', '3', '3', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '.', '.', '.', '.', '#', '.', '.'})

		for i := 0; i < b.N; i++ {
			_, _ = extractTokens(schematic)
		}
	})

	b.Run("gears extraction", func(b *testing.B) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '*', '.', '.', '.', '*', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '3', '5', '.', '.', '.', '3', '3', '.'})
		schematic = append(schematic, []byte{'.', '*', '.', '.', '.', '*', '.', '#', '.', '.'})

		for i := 0; i < b.N; i++ {
			_, _ = extractTokens(schematic)
		}
	})

	b.Run("value adjacency determination", func(b *testing.B) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '#'})
		schematic = append(schematic, []byte{'.', '.', '.', '.', '*', '3', '2', '.', '5', '6'})
		schematic = append(schematic, []byte{'.', '.', '5', '0', '1', '.', '1', '2', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '.', '.', '.', '/', '.', '.', '6'})

		for i := 0; i < b.N; i++ {
			isAdjacentToSymbols(schematicValue{"467", 0, 0}, schematic)
		}
	})

	b.Run("gear ratio determination", func(b *testing.B) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '1', '1', '4', '.', '*'})
		schematic = append(schematic, []byte{'.', '*', '.', '.', '*', '3', '2', '.', '5', '6'})
		schematic = append(schematic, []byte{'.', '.', '5', '0', '1', '.', '1', '2', '.', '*'})
		schematic = append(schematic, []byte{'*', '.', '.', '.', '.', '.', '/', '.', '.', '6'})

		for i := 0; i < b.N; i++ {
			getGearRatioOrZero(gear{0, 9}, schematic)
		}
	})

	b.Run("schematic value from coordinates extraction", func(b *testing.B) {
		var schematic [][]byte
		schematic = append(schematic, []byte{'4', '6', '7', '.', '.', '.', '.', '.', '.', '#'})
		schematic = append(schematic, []byte{'.', '.', '.', '.', '*', '3', '2', '.', '5', '6'})
		schematic = append(schematic, []byte{'.', '.', '5', '0', '1', '.', '1', '.', '.', '.'})
		schematic = append(schematic, []byte{'.', '.', '.', '.', '.', '.', '/', '.', '.', '6'})

		for i := 0; i < b.N; i++ {
			_, _ = getValueFromCoordinates(0, 0, schematic)
		}
	})

}
