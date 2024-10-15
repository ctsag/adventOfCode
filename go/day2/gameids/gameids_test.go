package gameids

import (
	"errors"
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

func TestTotalsCalculationShould(t *testing.T) {

	t.Run("calculate total of possible game ids", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = `Game 1: 7 blue, 9 red, 1 green; 8 green; 10 green, 5 blue, 3 red; 11 blue, 5 red, 1 green
					   Game 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green
					   Game 3: 11 blue, 3 red, 1 green; 15 red, 9 blue, 3 green; 11 blue, 4 red, 4 green; 1 red, 2 green, 14 blue; 18 blue, 4 green, 10 red
					   Game 4: 3 red, 7 blue; 3 blue, 2 red, 2 green; 2 green, 1 red, 1 blue; 3 green, 5 blue, 5 red; 7 blue, 1 green, 1 red; 2 green, 7 blue
					   Game 5: 1 blue, 2 red, 1 green; 6 blue, 3 green, 2 red; 2 blue`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _, _ := CalculateTotals(fileName, mockReader)
		expected := 1 + 4 + 5

		assert.Equal(
			t,
			expected,
			actual,
			"Did not calculate the game id total correctly",
		)
	})

	t.Run("calculate total minimum cubes total", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = `Game 1: 3 blue, 4 red; 1 red, 2 green, 6 blue; 2 green
                       Game 2: 1 blue, 2 green; 3 green, 4 blue, 1 red; 1 green, 1 blue
                       Game 3: 6 red, 1 blue, 3 green; 2 blue, 1 red, 2 green`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		_, actual, _ := CalculateTotals(fileName, mockReader)
		expected := 48 + 12 + 36

		assert.Equal(
			t,
			expected,
			actual,
			"Did not calculate the minimum cubes total correctly",
		)
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

	t.Run("fail when unable to parse game id", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = `Game 1: 7 blue, 9 red, 1 green; 8 green; 10 green, 5 blue, 3 red; 11 blue, 5 red, 1 green
	                   Round 2: 3 red, 7 blue; 3 blue, 2 red, 2 green; 2 green, 1 red, 1 blue; 3 green, 5 blue, 5 red; 7 blue, 1 green, 1 red; 2 green, 7 blue
			           Game 3: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		_, _, err := CalculateTotals(fileName, mockReader)
		expected := "unable to extract game id in line"

		assert.ErrorContains(
			t,
			err,
			expected,
			"Did not handle game id parsing failure properly",
		)
	})

	t.Run("fail when unable to parse color values", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = `Game 1: 7 blue, 9 red, 1 green; 8 green; 10 green, 5 blue, 3 red; 11 blue, 5 red, 1 green
	                   Game 2- 3 red, 7 blue; 3 blue, 2 red, 2 green; 2 green, 1 red, 1 blue; 3 green, 5 blue, 5 red; 7 blue, 1 green, 1 red; 2 green, 7 blue
			           Game 3: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		_, _, err := CalculateTotals(fileName, mockReader)
		expected := "unable to extract color values in line"

		assert.ErrorContains(
			t,
			err,
			expected,
			"Did not handle game id parsing failure properly",
		)
	})

}

func TestPossibleGameDeterminationShould(t *testing.T) {

	t.Run("return id when game is possible", func(t *testing.T) {
		const line = "Game 4: 3 red, 7 blue; 3 blue, 2 red, 2 green; 2 green, 1 red, 1 blue; 3 green, 5 blue, 5 red; 7 blue, 1 green, 1 red; 2 green, 7 blue"

		actual, _ := possibleGameIdOrZero(line)
		expected := 4

		assert.Equal(
			t,
			expected,
			actual,
			"Did not extract id when game is possible",
		)
	})

	t.Run("return 0 when game is not possible", func(t *testing.T) {
		const line = "Game 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		actual, _ := possibleGameIdOrZero(line)
		expected := 0

		assert.Equal(
			t,
			expected,
			actual,
			"Did not return 0 when game is not possible",
		)
	})

	t.Run("fail when unable to extract color values", func(t *testing.T) {
		const line = "Game 2- 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		_, err := possibleGameIdOrZero(line)
		expected := "unable to extract color values in line"

		assert.ErrorContains(
			t,
			err,
			expected,
			"Did not handle game id parsing failure properly",
		)
	})

	t.Run("fail when unable to extract game id", func(t *testing.T) {
		const line = "Round 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		_, err := possibleGameIdOrZero(line)
		expected := "unable to extract game id in line"

		assert.ErrorContains(
			t,
			err,
			expected,
			"Did not handle game id parsing failure properly",
		)
	})

}

func TestMaxColorValuesExtractionShould(t *testing.T) {

	t.Run("extract maximum values per color", func(t *testing.T) {
		const line = "Game 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		actual, _ := extractMaxColorValues(line)
		expected := map[string]int{"red": 6, "green": 7, "blue": 20}

		assert.Equal(
			t,
			expected,
			actual,
			"Did not extract max color values correctly",
		)
	})

	t.Run("fail when unable to split the line provided", func(t *testing.T) {
		const line = "7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		_, err := extractMaxColorValues(line)
		expected := "unable to parse line"

		assert.ErrorContains(
			t,
			err,
			expected,
			"Did not extract max color values correctly",
		)
	})

}

func TestGameIdExtractionShould(t *testing.T) {

	t.Run("extract game id", func(t *testing.T) {
		const line = "Game 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		actual, _ := extractGameId(line)
		expected := 2

		assert.Equal(
			t,
			expected,
			actual,
			"Did not extract game id correctly",
		)
	})

	t.Run("fail to extract game id when not provided", func(t *testing.T) {
		const line = "Round 1: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		_, err := extractGameId(line)
		expected := "game id not found"

		assert.EqualError(
			t,
			err,
			expected,
			"Did not fail to extract game id properly",
		)
	})

}

func TestMinPossibleCubesDeterminationShould(t *testing.T) {

	t.Run("return minimum possible cubes", func(t *testing.T) {
		const line = "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red"

		actual, _ := minimumPossibleCubes(line)
		expected := 20 * 13 * 6

		assert.Equal(
			t,
			expected,
			actual,
			"Did not calculate minimum possible cubes correctly",
		)
	})

	t.Run("fail when unable to extract color values", func(t *testing.T) {
		const line = "Game 2- 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		_, err := minimumPossibleCubes(line)
		expected := "unable to extract color values in line"

		assert.ErrorContains(
			t,
			err,
			expected,
			"Did not handle game id parsing failure properly",
		)
	})

}

func BenchmarkTotalsCalculation(b *testing.B) {

	b.Run("totals calculation", func(b *testing.B) {
		const fileName = "test_input.txt"
		const lines = `Game 1: 7 blue, 9 red, 1 green; 8 green; 10 green, 5 blue, 3 red; 11 blue, 5 red, 1 green
					   Game 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green
					   Game 3: 11 blue, 3 red, 1 green; 15 red, 9 blue, 3 green; 11 blue, 4 red, 4 green; 1 red, 2 green, 14 blue; 18 blue, 4 green, 10 red
					   Game 4: 3 red, 7 blue; 3 blue, 2 red, 2 green; 2 green, 1 red, 1 blue; 3 green, 5 blue, 5 red; 7 blue, 1 green, 1 red; 2 green, 7 blue
					   Game 5: 1 blue, 2 red, 1 green; 6 blue, 3 green, 2 red; 2 blue`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		for i := 0; i < b.N; i++ {
			_, _, _ = CalculateTotals(fileName, mockReader)
		}
	})

	b.Run("possible game determination", func(b *testing.B) {
		const line = "Game 4: 3 red, 7 blue; 3 blue, 2 red, 2 green; 2 green, 1 red, 1 blue; 3 green, 5 blue, 5 red; 7 blue, 1 green, 1 red; 2 green, 7 blue"

		for i := 0; i < b.N; i++ {
			_, _ = possibleGameIdOrZero(line)
		}
	})

	b.Run("max color values extraction", func(b *testing.B) {
		const line = "Game 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		for i := 0; i < b.N; i++ {
			_, _ = extractMaxColorValues(line)
		}
	})

	b.Run("game id extraction", func(b *testing.B) {
		const line = "Game 2: 7 green, 3 blue; 20 blue, 4 green; 6 red, 13 blue, 2 green"

		for i := 0; i < b.N; i++ {
			_, _ = extractGameId(line)
		}
	})

	b.Run("min possible cubes determination", func(b *testing.B) {
		const line = "Game 3: 8 green, 6 blue, 20 red; 5 blue, 4 red, 13 green; 5 green, 1 red"

		for i := 0; i < b.N; i++ {
			_, _ = minimumPossibleCubes(line)
		}
	})

}
