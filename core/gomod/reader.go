package gomod

import "os"

type Reader interface {
	Read(path string) ([]byte, error)
}

type FileReader struct{}

func NewReader() Reader {
	return &FileReader{}
}

func (r *FileReader) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}
