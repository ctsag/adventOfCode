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

func TestFileOpsShould(t *testing.T) {

	t.Run("open file", func(t *testing.T) {
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
	})

	t.Run("fail when unable to open file", func(t *testing.T) {
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
	})

	t.Run("suppress failure to close file", func(t *testing.T) {
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
	})

}

func TestConcreteFileReaderShould(t *testing.T) {

	t.Run("open file", func(t *testing.T) {
		const fileName = "../testdata/test_input.txt"

		fileReader := FileReader{}

		_, err := fileReader.Open(fileName)

		assert.Nil(
			t,
			err,
			"Failed to open file",
		)
	})

}

func BenchmarkFileOps(b *testing.B) {

	b.Run("mock reader", func(b *testing.B) {
		const fileName = "test_input.txt"
		const contents = "test contents"

		mockReader := new(MockFileReader)
		mockFile := new(MockReadCloser)

		mockFile.On("Close").Return(nil)
		mockReader.On("Open", fileName).Return(io.NopCloser(strings.NewReader(contents)), nil)

		for i := 0; i < b.N; i++ {
			_, _ = OpenFile(fileName, mockReader)
		}
	})

	b.Run("concrete reader", func(b *testing.B) {
		const fileName = "../testdata/test_input.txt"

		fileReader := FileReader{}

		for i := 0; i < b.N; i++ {
			_, _ = fileReader.Open(fileName)
		}
	})

}
