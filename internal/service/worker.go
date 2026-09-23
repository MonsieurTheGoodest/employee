package service

import (
	"context"
	"log"
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
				log.Printf("worker ERR: %s", err.Error())
			}

			if len(checkedEmployees) > 0 {
				for _, emp := range checkedEmployees {
					if err := s.db.ChangeStatus(ctx, emp); err != nil {
						log.Printf("worker ERR: changing status of %v ERR: %s", emp, err.Error())
					}
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
