package daemon

import (
	"context"
	"fmt"
	"synctl/internal/application/interfaces"
	"synctl/internal/application/services/daemon/reconciler"
	"time"
)

type Controller struct {
	repo       interfaces.Repository
	reconciler *reconciler.ReconcileService
	interval   time.Duration
	logger     interfaces.Logger
}

func NewController(
	repo interfaces.Repository,
	reconciler *reconciler.ReconcileService,
	interval time.Duration,
	logger interfaces.Logger,
) *Controller {
	return &Controller{
		repo:       repo,
		reconciler: reconciler,
		interval:   interval,
		logger:     logger,
	}
}

func (c *Controller) Start(ctx context.Context) {
	c.logger.Info("Controller - reconciling state...")
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			state, err := c.repo.Load()
			if err != nil {
				continue
			}

			if err := c.reconciler.Reconcile(state); err != nil {
				c.logger.Error(fmt.Sprintf("Error reconcile: %s", err))
				continue
			}
		case <-ctx.Done():
			return
		}
	}
}
