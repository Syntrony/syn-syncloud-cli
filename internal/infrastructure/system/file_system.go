package system

import "os"

type FileSystem struct{}

func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

func (f *FileSystem) WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0664)
}

func (f *FileSystem) ExistsFile(path string) (bool, error) {
	_, err := os.Stat(path)

	if err != nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
