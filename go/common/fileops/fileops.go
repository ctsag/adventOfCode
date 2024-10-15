package fileops

import (
	"io"
	"os"
)

type ReadableFile interface {
	Open(filePath string) (io.ReadCloser, error)
}

type FileReader struct{}

func (f *FileReader) Open(filePath string) (io.ReadCloser, error) {
	return os.Open(filePath)
}

func OpenFile(path string, reader ReadableFile) (io.ReadCloser, error) {
	file, err := reader.Open(path)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func CloseFile(file io.ReadCloser) error {
	if file != nil {
		_ = file.Close()
	}

	return nil
}
