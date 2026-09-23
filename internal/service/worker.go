package service

import (
	"context"
	"fmt"
	"time"
)

func (s *Service) Checking(
	ctx context.Context,
	checkInterval time.Duration,
	pendingTimeInSeconds int,
) error {

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			checkedEmployees, err := s.db.CheckPendingEmployees(ctx, pendingTimeInSeconds)

			if err != nil {
				return fmt.Errorf("worker ERR: %w", err)
			}

			if len(checkedEmployees) > 0 {
				for _, emp := range checkedEmployees {
					if err := s.db.ChangeStatus(ctx, emp); err != nil {
						return fmt.Errorf("worker ERR: changing status of %v ERR: %w", emp, err)
					}
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
