package repository

import (
	"synctl/internal/infrastructure/state"
)

func StateRepository() state.Repository {
	return &state.StateRepository{
		Path: "helpers/state.json",
	}
}
