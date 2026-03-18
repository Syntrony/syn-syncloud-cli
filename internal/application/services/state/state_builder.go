package state

import (
	"synctl/internal/domain"
	"time"
)

type StateBuilder struct{}

func NewStateBuilder() *StateBuilder {
	return &StateBuilder{}
}

func (b *StateBuilder) Build(snapshot *domain.Snapshot, nodes []domain.Node) *domain.State {
	cluster := &domain.Cluster{
		Id:        "cluster-01",
		Name:      "syncloud-local",
		Mode:      "single",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	return &domain.State{
		Version:   "0.1",
		Cluster:   cluster,
		Nodes:     nodes,
		Resources: snapshot.Resources,
		Dns:       &domain.Dns{Records: snapshot.Records},
	}
}
