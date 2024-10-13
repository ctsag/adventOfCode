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
	const lines = "7pqrstsixteen\neightwothree\nzoneight234"

	mockReader := new(MockFileReader)
	mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(lines)), nil)

	actual, _ := CalculateTotal(fileName, mockReader)
	expected := 76 + 83 + 14

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
		"Did not fail when unable to read file",
	)
}

func TestCombinesWhenOnlyTwoDigitsAreProvided(t *testing.T) {
	const line = "aXonebcdefghi9j"

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
	const line = "oid7afbk3ceeightao"

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
	const line = "gsFgsixasboeomNa"

	actual := combineFirstAndLastDigit(line)
	expected := 66

	assert.Equal(
		t,
		expected,
		actual,
		"Did not combine the first and last digit correctly",
	)
}

func TestReplacesWordsWithDigits(t *testing.T) {
	const input = "5five_sixseven8one1twozthreefoureight9nine0eightwozero"

	actual := replaceWordsWithDigits(input)
	expected := "5f5ive_s6ixs7even8o1ne1t2wozt3hreef4oure8ight9n9ine0e8ight2wozero"

	assert.Equal(
		t,
		expected,
		actual,
		"Did not properly replace words with digits",
	)
}

func TestRetrievesIndexOfItemInArray(t *testing.T) {
	arr := []string{"one", "two", "three"}
	val := "two"

	actual, _ := arrayIndexOf(arr, val)
	expected := 1

	assert.Equal(
		t,
		expected,
		actual,
		"Did not retrieve index of item in array",
	)
}

func TestFailsWhenItemNotInArray(t *testing.T) {
	arr := []string{"one", "two", "three"}
	val := "zero"

	_, err := arrayIndexOf(arr, val)
	expected := "index not found"

	assert.EqualError(
		t,
		err,
		expected,
		"Did not fail when unable to retrieve index of item in array",
	)
}

func TestReplacesInANonDestructiveWaY(t *testing.T) {
	val := "five"
	num := 5

	expected := "f5ive"
	actual := nonDestructiveReplace(val, num)

	assert.Equal(
		t,
		expected,
		actual,
		"Did not replace properly",
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
