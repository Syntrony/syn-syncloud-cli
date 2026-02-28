package repository

import (
	"synctl/internal/domain/interfaces"
	"synctl/internal/state"
)

func StateRepository() interfaces.Repository {
	return &state.StateRepository{
		Path: "helpers/state.json",
	}
}
