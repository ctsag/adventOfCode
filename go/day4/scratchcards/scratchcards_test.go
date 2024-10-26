package scratchcards

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

	t.Run("calculate total of scratchcards", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = "Card 1: 41 48  | 83 41 6 31 17\nCard 2: 13 32 | 61 30 68 82 17\nCard 3: 1 21 15 | 15 21 63 1 16"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _, _ := CalculateTotals(fileName, mockReader)
		expected := 1 + 0 + 4

		assert.Equal(t, expected, actual, "Did not calculate the scratchcards total correctly")
	})

	t.Run("calculate count of bonus scratchcards", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = "Card 1: 41 48  | 43 41 6 31 17\nCard 2: 13 32 | 61 30 68 82 17\nCard 3: 1 21 15 | 15 21 63 1 16"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _, _ := CalculateTotals(fileName, mockReader)
		expected := 1 + 2 + 2

		assert.Equal(t, expected, actual, "Did not calculate the count of bonus scratchcards correctly")
	})

	t.Run("fails when unable to read file", func(t *testing.T) {
		const fileName = "test_input.txt"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader("")), errors.New("file open error"))

		_, _, err := CalculateTotals(fileName, mockReader)
		expected := "file open error"

		assert.EqualError(t, err, expected, "Did not fail when unable to read file")
	})

}

func TestScratchcardExtractionShould(t *testing.T) {

	t.Run("extract scratchcards", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = "Card 1: 41 48  | 83 41 6 31 17\nCard 2: 13 32 | 61 30 68 82 17"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _ := extractScratchcards(fileName, mockReader)

		var expected []scratchcard
		expected = append(expected, scratchcard{1, []int{41, 48}, []int{83, 41, 6, 31, 17}})
		expected = append(expected, scratchcard{2, []int{13, 32}, []int{61, 30, 68, 82, 17}})

		assert.Equal(t, expected, actual, "Did not extract scratchcards correctly")
	})

	t.Run("extract card from line", func(t *testing.T) {
		const line = "Card 1: 41 48  | 83 41 6 31 17"

		actual := cardFromLine(line)
		expected := scratchcard{1, []int{41, 48}, []int{83, 41, 6, 31, 17}}

		assert.Equal(t, expected, actual, "Did not extract scratchcard from line")
	})

	t.Run("convert to int array", func(t *testing.T) {
		const line = " 83 41  6 31   17"

		actual := toIntArray(line)
		expected := []int{83, 41, 6, 31, 17}

		assert.Equal(t, expected, actual, "Did not convert to int array")
	})

}

func TestScoreDeterminationShould(t *testing.T) {

	tests := []struct {
		card  scratchcard
		score int
	}{
		{scratchcard{1, []int{41, 48}, []int{58, 61, 3}}, 0},
		{scratchcard{2, []int{41, 48}, []int{83, 41, 6}}, 1},
		{scratchcard{3, []int{41, 48}, []int{48, 41, 6}}, 2},
		{scratchcard{4, []int{41, 48, 56}, []int{48, 41, 56}}, 4},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("determine card score for %v", test.card), func(t *testing.T) {
			actual := determineScore(test.card)
			expected := test.score

			assert.Equal(t, expected, actual, "Did not determine score correctly")
		})
	}

}

func TestBonusCountDeterminationShould(t *testing.T) {

	t.Run("determine count of bonus cards", func(t *testing.T) {
		scratchcards := []scratchcard{
			{1, []int{41, 48, 83, 86, 17}, []int{83, 86, 6, 31, 17, 9, 48, 53}},
			{2, []int{13, 32, 20, 16, 61}, []int{61, 30, 68, 82, 17, 32, 24, 19}},
			{3, []int{1, 21, 53, 59, 44}, []int{69, 82, 63, 72, 16, 21, 14, 1}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{6, []int{31, 18, 13, 56, 72}, []int{74, 77, 10, 23, 35, 67, 36, 11}},
		}

		expected := []scratchcard{
			{1, []int{41, 48, 83, 86, 17}, []int{83, 86, 6, 31, 17, 9, 48, 53}},
			{2, []int{13, 32, 20, 16, 61}, []int{61, 30, 68, 82, 17, 32, 24, 19}},
			{3, []int{1, 21, 53, 59, 44}, []int{69, 82, 63, 72, 16, 21, 14, 1}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{6, []int{31, 18, 13, 56, 72}, []int{74, 77, 10, 23, 35, 67, 36, 11}},
			{2, []int{13, 32, 20, 16, 61}, []int{61, 30, 68, 82, 17, 32, 24, 19}},
			{3, []int{1, 21, 53, 59, 44}, []int{69, 82, 63, 72, 16, 21, 14, 1}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{3, []int{1, 21, 53, 59, 44}, []int{69, 82, 63, 72, 16, 21, 14, 1}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{3, []int{1, 21, 53, 59, 44}, []int{69, 82, 63, 72, 16, 21, 14, 1}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
		}

		actual := addBonusCards(scratchcards)

		assert.Equal(t, expected, actual, "Did not determine bonus cards correctly")
	})

	tests := []struct {
		card    scratchcard
		matches int
	}{
		{scratchcard{1, []int{41, 48, 64}, []int{48, 41, 64}}, 3},
		{scratchcard{1, []int{41, 48, 64}, []int{48, 41, 61}}, 2},
		{scratchcard{1, []int{41, 48, 64}, []int{48, 41, 41}}, 2},
		{scratchcard{1, []int{41, 48, 64}, []int{50, 60, 61}}, 0},
		{scratchcard{1, []int{41, 48, 64}, []int{48, 51, 61}}, 1},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("determine winning matches for scratchcard %v", test.card), func(t *testing.T) {
			actual := countMatches(test.card)
			expected := test.matches

			assert.Equal(t, expected, actual, "Did not determine winning matches correctly")
		})
	}

}

func BenchmarkTotalsCalculation(b *testing.B) {

	b.Run("totals calculation", func(b *testing.B) {
		const fileName = "test_input.txt"
		const lines = "Card 1: 41 48  | 83 41 6 31 17\nCard 2: 13 32 | 61 30 68 82 17\nCard 3: 1 21 15 | 15 21 63 1 16"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		for i := 0; i < b.N; i++ {
			_, _, _ = CalculateTotals(fileName, mockReader)
		}
	})

	b.Run("scratchcards extraction", func(b *testing.B) {
		const fileName = "test_input.txt"
		const lines = "Card 1: 41 48  | 83 41 6 31 17\nCard 2: 13 32 | 61 30 68 82 17"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		for i := 0; i < b.N; i++ {
			_, _ = extractScratchcards(fileName, mockReader)
		}
	})

	b.Run("scratchcards extraction from line", func(b *testing.B) {
		const line = "Card 1: 41 48  | 83 41 6 31 17"

		for i := 0; i < b.N; i++ {
			_ = cardFromLine(line)
		}
	})

	b.Run("conversion to int array", func(b *testing.B) {
		const line = " 83 41  6 31   17"

		for i := 0; i < b.N; i++ {
			_ = toIntArray(line)
		}
	})

	b.Run("score determination", func(b *testing.B) {
		card := scratchcard{1, []int{41, 48}, []int{58, 61, 3}}

		for i := 0; i < b.N; i++ {
			_ = determineScore(card)
		}
	})

	b.Run("bonus card count", func(b *testing.B) {
		scratchcards := []scratchcard{
			{1, []int{41, 48, 83, 86, 17}, []int{83, 86, 6, 31, 17, 9, 48, 53}},
			{2, []int{13, 32, 20, 16, 61}, []int{61, 30, 68, 82, 17, 32, 24, 19}},
			{3, []int{1, 21, 53, 59, 44}, []int{69, 82, 63, 72, 16, 21, 14, 1}},
			{4, []int{41, 92, 73, 84, 69}, []int{59, 84, 76, 51, 58, 5, 54, 83}},
			{5, []int{87, 83, 26, 28, 32}, []int{88, 30, 70, 12, 93, 22, 82, 36}},
			{6, []int{31, 18, 13, 56, 72}, []int{74, 77, 10, 23, 35, 67, 36, 11}},
		}

		for i := 0; i < b.N; i++ {
			_ = addBonusCards(scratchcards)
		}
	})

	b.Run("matching numbers count", func(b *testing.B) {
		card := scratchcard{1, []int{41, 48, 64}, []int{48, 41, 64}}

		for i := 0; i < b.N; i++ {
			_ = countMatches(card)
		}
	})

}
