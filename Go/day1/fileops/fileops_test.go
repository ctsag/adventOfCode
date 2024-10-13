package fileops

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"strings"
	"testing"
)

type MockReadCloser struct {
	mock.Mock
	io.Reader
}

func (mockCloser *MockReadCloser) Close() error {
	return mockCloser.Called().Error(0)
}

type MockFileReader struct {
	mock.Mock
}

func (mockReader *MockFileReader) Open(path string) (io.ReadCloser, error) {
	args := mockReader.Called(path)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func TestOpensFile(t *testing.T) {
	const fileName = "test_input.txt"
	const contents = "test contents"

	mockReader := new(MockFileReader)
	mockFile := new(MockReadCloser)

	mockFile.On("Close").Return(nil)
	mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(contents)), nil)

	reader, err := OpenFile(fileName, mockReader)

	assert.Nil(
		t,
		err,
		"Did not open file",
	)

	actual, _ := io.ReadAll(reader)
	expected := contents

	assert.Equal(
		t,
		string(actual),
		expected,
		"Did not return expected contents",
	)
}

func TestHandlesFileOpenFailure(t *testing.T) {
	const fileName = "test_input.txt"

	mockReader := new(MockFileReader)
	mockFile := new(MockReadCloser)

	mockFile.On("Close").Return(nil)
	mockReader.On("Open", fileName).Return(io.NopCloser(nil), errors.New("file open error"))

	_, err := OpenFile(fileName, mockReader)
	expected := "file open error"

	assert.EqualError(
		t,
		err,
		expected,
		"Did not return expected error",
	)
}

func TestSuppressesFileCloseFailure(t *testing.T) {
	const fileName = "test_input.txt"

	mockReader := new(MockFileReader)
	mockFile := new(MockReadCloser)

	mockFile.On("Close").Return(errors.New("file close error"))
	mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader("")), nil)

	file, _ := OpenFile(fileName, mockReader)
	err := CloseFile(file)

	assert.Nil(
		t,
		err,
		"Did not suppress file close error",
	)
}
