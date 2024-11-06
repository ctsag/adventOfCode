package almanac

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

func TestLowestLocationDeterminationShould(t *testing.T) {

	t.Run("determine lowest possible locations", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _ := DetermineLowestLocation(fileName, mockReader)
		expected := 35

		assert.Equal(t, expected, actual, "Did not determine the lowest possible locations correctly")
	})

	t.Run("fails when unable to read file", func(t *testing.T) {
		const fileName = "test_input.txt"

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader("")), errors.New("file open error"))

		_, err := DetermineLowestLocation(fileName, mockReader)
		expected := "file open error"

		assert.EqualError(t, err, expected, "Did not fail when unable to read file")
	})

}

func TestAlmanacExtractionShould(t *testing.T) {

	t.Run("extract the almanac", func(t *testing.T) {
		const fileName = "test_input.txt"
		const lines = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		actual, _ := extractAlmanac(fileName, mockReader)
		expected := almanac{
			[]int{79, 14, 55, 13},
			[]rangeMap{
				{50, 98, 2},
				{52, 50, 48},
			},
			[]rangeMap{
				{0, 15, 37},
				{37, 52, 2},
				{39, 0, 15},
			},
			[]rangeMap{
				{49, 53, 8},
				{0, 11, 42},
				{42, 0, 7},
				{57, 7, 4},
			},
			[]rangeMap{
				{88, 18, 7},
				{18, 25, 70},
			},
			[]rangeMap{
				{45, 77, 23},
				{81, 45, 19},
				{68, 64, 13},
			},
			[]rangeMap{
				{0, 69, 1},
				{1, 0, 69},
			},
			[]rangeMap{
				{60, 56, 37},
				{56, 93, 4},
			},
		}

		assert.Equal(t, expected, actual, "Did not extract seeds maps correctly")
	})

	t.Run("convert string to int array", func(t *testing.T) {
		line := "    32       45 51     "

		actual := toIntArray(line)
		expected := []int{32, 45, 51}

		assert.Equal(t, expected, actual, "Did not convert string to int array")
	})

}

func TestSeedExtractionShould(t *testing.T) {

	t.Run("extract seeds", func(t *testing.T) {
		id := 79
		alm := almanac{
			[]int{79, 14, 55, 13},
			[]rangeMap{
				{52, 50, 48},
			},
			[]rangeMap{
				{0, 15, 37},
			},
			[]rangeMap{
				{49, 53, 8},
			},
			[]rangeMap{
				{18, 25, 70},
			},
			[]rangeMap{
				{68, 64, 13},
			},
			[]rangeMap{
				{0, 69, 1},
			},
			[]rangeMap{
				{60, 56, 37},
				{56, 93, 4},
			},
		}

		actual := extractSeed(alm, id)
		expected := seed{
			79,
			81,
			81,
			81,
			74,
			78,
			78,
			82,
		}

		assert.Equal(t, expected, actual, "Did not extract seeds maps correctly")
	})

}

func TestNextValueExtractionShould(t *testing.T) {

	tests := []struct {
		rmap     []rangeMap
		id       int
		expected int
	}{
		{
			[]rangeMap{
				{50, 98, 2},
				{52, 50, 48},
			},
			79,
			81,
		},
		{
			[]rangeMap{
				{45, 77, 23},
				{81, 45, 19},
				{68, 64, 13},
			},
			74,
			78,
		},
		{
			[]rangeMap{
				{49, 53, 8},
				{0, 11, 42},
				{42, 0, 7},
				{57, 7, 4},
			},
			53,
			49,
		},
		{
			[]rangeMap{
				{0, 69, 1},
			},
			78,
			78,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("extracts next value for %d from %v", test.id, test.rmap), func(t *testing.T) {
			actual := nextValue(test.rmap, test.id)

			assert.Equal(t, test.expected, actual, "Did not extract seeds correctly")
		})
	}

}

func BenchmarkLowestLocationCalculation(b *testing.B) {

	b.Run("lowest location calculation", func(b *testing.B) {
		const fileName = "test_input.txt"
		const lines = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		for i := 0; i < b.N; i++ {
			_, _ = DetermineLowestLocation(fileName, mockReader)
		}
	})

	b.Run("almanac extraction", func(b *testing.B) {
		const fileName = "test_input.txt"
		const lines = `seeds: 79 14 55 13

seed-to-soil map:
50 98 2
52 50 48

soil-to-fertilizer map:
0 15 37
37 52 2
39 0 15

fertilizer-to-water map:
49 53 8
0 11 42
42 0 7
57 7 4

water-to-light map:
88 18 7
18 25 70

light-to-temperature map:
45 77 23
81 45 19
68 64 13

temperature-to-humidity map:
0 69 1
1 0 69

humidity-to-location map:
60 56 37
56 93 4`

		mockReader := new(MockFileReader)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

		for i := 0; i < b.N; i++ {
			_, _ = extractAlmanac(fileName, mockReader)
		}
	})

	b.Run("int array from string extraction", func(b *testing.B) {
		line := "    32       45 51     "

		for i := 0; i < b.N; i++ {
			_ = toIntArray(line)
		}
	})

	b.Run("seed extraction", func(b *testing.B) {
		id := 79
		alm := almanac{
			[]int{79, 14, 55, 13},
			[]rangeMap{
				{52, 50, 48},
			},
			[]rangeMap{
				{0, 15, 37},
			},
			[]rangeMap{
				{49, 53, 8},
			},
			[]rangeMap{
				{18, 25, 70},
			},
			[]rangeMap{
				{68, 64, 13},
			},
			[]rangeMap{
				{0, 69, 1},
			},
			[]rangeMap{
				{60, 56, 37},
				{56, 93, 4},
			},
		}

		for i := 0; i < b.N; i++ {
			_ = extractSeed(alm, id)
		}
	})

	b.Run("next value extraction", func(b *testing.B) {
		rmap := []rangeMap{{50, 98, 2}}
		id := 79

		for i := 0; i < b.N; i++ {
			_ = nextValue(rmap, id)
		}
	})
}
