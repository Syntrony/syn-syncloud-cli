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

func (f *StateRepository) Upsert(desired []*domain.Resource) ([]domain.Resource, error) {
	state, err := f.Load()

	if err != nil {
		return nil, err
	}

	merged := map[string]domain.Resource{}

	for _, res := range state.Resources {
		merged[res.Key()] = res
	}

	for _, res := range desired {
		key := res.Key()

		if existing, ok := merged[key]; ok {
			res.Id = existing.Id
			res.Name = existing.Name
			res.Kind = existing.Kind
			res.Runtime = existing.Runtime
			res.NodeId = existing.NodeId
			res.Spec = existing.Spec
			res.Status = existing.Status
			res.Ownership = existing.Ownership
			res.CreatedAt = existing.CreatedAt
			res.UpdatedAt = existing.UpdatedAt
		}
		merged[key] = *res
	}

	var snapshot []domain.Resource
	for _, res := range merged {
		snapshot = append(snapshot, res)
	}

	return snapshot, nil
}

func (f *StateRepository) Remove(desired []*domain.Resource) ([]domain.Resource, error) {
	state, err := f.Load()

	if err != nil {
		return nil, err
	}

	removeSet := make(map[string]struct{})

	for _, res := range desired {
		removeSet[res.Key()] = struct{}{}
	}

	var snapshot []domain.Resource

	for _, existing := range state.Resources {
		if _, shouldDelete := removeSet[existing.Key()]; shouldDelete {
			continue
		}
		snapshot = append(snapshot, existing)
	}

	return snapshot, nil
}
