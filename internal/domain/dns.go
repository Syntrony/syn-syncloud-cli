package domain

import resource "synctl/internal/domain/resource"

type Dns struct {
	Records []resource.Record `json:"records,omitempty"`
}
