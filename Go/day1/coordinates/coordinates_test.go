package coordinates

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

func TestCalculatesTotal(t *testing.T) {
	const fileName = "test_input.txt"
	const lines = "aX1bcdefghi9j\ngsFg6asboeomNa\noid7afbk3ce8ao"

	mockReader := new(MockFileReader)
	mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

	actual, _ := CalculateTotal(fileName, mockReader)
	expected := 19 + 66 + 78

	assert.Equal(
		t,
		expected,
		actual,
		"Did not calculate the total correctly",
	)
}

func TestFailsWhenUnableToReadFile(t *testing.T) {
	const fileName = "test_input.txt"

	mockReader := new(MockFileReader)
	mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader("")), errors.New("file open error"))

	_, err := CalculateTotal(fileName, mockReader)
	expected := "file open error"

	assert.EqualError(
		t,
		err,
		expected,
	)
}

func TestCombinesWhenOnlyTwoDigitsAreProvided(t *testing.T) {
	const line = "aX1bcdefghi9j"

	actual := combineFirstAndLastDigit(line)
	expected := 19

	assert.Equal(
		t,
		expected,
		actual,
		"Did not combine the first and last digit correctly",
	)
}

func TestCombineWhenMoreThanTwoDigitsAreProvided(t *testing.T) {
	const line = "oid7afbk3ce8ao"

	actual := combineFirstAndLastDigit(line)
	expected := 78

	assert.Equal(
		t,
		expected,
		actual,
		"Did not combine the first and last digit correctly",
	)
}

func TestCombinesWhenOnlyOneDigitIsProvided(t *testing.T) {
	const line = "gsFg6asboeomNa"

	actual := combineFirstAndLastDigit(line)
	expected := 66

	assert.Equal(
		t,
		expected,
		actual,
		"Did not combine the first and last digit correctly",
	)
}

func TestCanIdentifyDigit(t *testing.T) {
	const character rune = '8'

	actual := isDigitOtherwiseZero(uint8(character))
	expected := 8

	assert.Equal(
		t,
		expected,
		actual,
		"Did not identify a digit",
	)
}

func TestCanIdentifyNonDigit(t *testing.T) {
	const character rune = 'x'

	actual := isDigitOtherwiseZero(uint8(character))
	expected := 0

	assert.Equal(
		t,
		expected,
		actual,
		"Did not identify a non-digit",
	)
}
