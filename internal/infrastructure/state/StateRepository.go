package state

import (
	"encoding/json"
	"os"
	"synctl/internal/domain"
)

type StateRepository struct {
	Path string
}

func (f *StateRepository) Exists() (bool, error) {
	_, err := os.Stat(f.Path)

	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}

func (f *StateRepository) Load() (*domain.State, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}

	var state domain.State
	err = json.Unmarshal(data, &state)

	if err != nil {
		return nil, err
	}
	return &state, nil
}
