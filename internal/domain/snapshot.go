package domain

import resource "synctl/internal/domain/resource"

type Snapshot struct {
	Resources []Resource
	Records   []resource.Record
}
