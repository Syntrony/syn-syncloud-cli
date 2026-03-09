package interfaces

import "synctl/internal/domain/dto"

type PrivilegedRunner interface {
	Run(cmd string, args ...string) (*dto.CommandResult, error)
	// WriteFile(path string, content string) (*dto.CommandResult, error)
}
