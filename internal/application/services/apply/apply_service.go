package apply

import (
	"synctl/internal/application/interfaces"
	"synctl/internal/application/interfaces/mutation"
	mutationpkg "synctl/internal/application/services/mutation"
)

type ApplyService struct {
	logger   interfaces.Logger
	mutation *mutationpkg.MutationService
}

func NewApplyService(logger interfaces.Logger, mutation *mutationpkg.MutationService) *ApplyService {
	return &ApplyService{
		logger:   logger,
		mutation: mutation,
	}
}

func (s *ApplyService) Apply(file string, mut mutation.StateMutation) error {
	s.logger.Info("Applying Resource...")
	return s.mutation.Execute(file, mut)
}
