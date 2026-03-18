package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"
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

func (f *StateRepository) Save(st *domain.State) error {
	os.MkdirAll(filepath.Dir(f.Path), 0755)

	data, err := json.MarshalIndent(st, "", "  ")

	if err != nil {
		return err
	}

	return os.WriteFile(f.Path, data, 0644)
}

func (f *StateRepository) Upsert(resources []*domain.Resource) (*domain.Snapshot, error) {
	snapshot := &domain.Snapshot{
		Resources: make([]domain.Resource, 0, len(resources)),
	}

	for _, r := range resources {
		if r != nil {
			snapshot.Resources = append(snapshot.Resources, *r)
		}
	}

	return snapshot, nil

}
