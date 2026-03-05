package install

import "synctl/internal/domain/dto"

type Detector interface {
	Detect() (*dto.Status, error)
}
