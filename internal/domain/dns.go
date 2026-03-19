package domain

import resource "synctl/internal/domain/Resource"

type Dns struct {
	Records []resource.Record `json:"records,omitempty"`
}
