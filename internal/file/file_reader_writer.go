package file

import "os"

type FileReaderWriter struct{}

func (f *FileReaderWriter) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (f *FileReaderWriter) WriteFile(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func NewFileReaderWriter() *FileReaderWriter {
	return &FileReaderWriter{}
}
