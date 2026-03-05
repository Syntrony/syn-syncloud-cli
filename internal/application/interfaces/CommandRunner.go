package interfaces

import (
	dto "synctl/internal/domain/dto"
)

type CommandRunner interface {
	Run(name string, args ...string) (*dto.CommandResult, error)
}
