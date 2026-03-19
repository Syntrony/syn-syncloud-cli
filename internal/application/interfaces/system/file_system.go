package system

type FileSystem interface {
	WriteFile(path string, content string) error
	ExistsFile(path string) (bool, error)
	// BackupFile(path string) (string, error)
}
