package repository

import (
	"context"
	"fmt"
	"time"
)

func Checking(
	db *DataBase,
	checkInterval time.Duration,
	pendingTimeInSeconds int,
) error {

	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	defer cancel()

	for {
		select {
		case <-ticker.C:
			query := `
				SELECT employee_id
				FROM pending_employees 
				WHERE creating_time < NOW() - $1 * INTERVAL '1 second'
			`

			rows, err := db.Pool.Query(ctx, query, pendingTimeInSeconds)
			if err != nil {
				return fmt.Errorf("worker ERR: selecting employees ERR: %w", err)
			}
			defer rows.Close()

			var checkedEmployees []int
			for rows.Next() {
				var id int

				err := rows.Scan(&id)

				if err != nil {
					return fmt.Errorf("worker ERR: scanning employees ERR: %w", err)
				}

				checkedEmployees = append(checkedEmployees, id)
			}

			if len(checkedEmployees) > 0 {
				for _, emp := range checkedEmployees {
					if err := db.changeStatus(ctx, emp); err != nil {
						return fmt.Errorf("worker ERR: changing status of %v ERR: %w", emp, err)
					}
				}
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
