package domain

import resource "synctl/internal/domain/Resource"

type Snapshot struct {
	Resources []Resource
	Records   []resource.Record
}
