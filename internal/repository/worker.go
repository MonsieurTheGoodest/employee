package repository

import (
	"context"
	"employee/internal"
	"fmt"
	"time"
)

func Checking(db *DataBase) error {
	ctx, cancel := context.WithCancel(context.Background())
	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()
	defer cancel()

	for {
		select {
		case <-ticker.C:
			query := `
				SELECT first_name, last_name
				FROM pending_employees 
				WHERE creating_time < NOW() - $1 * INTERVAL '1 second'
			`
			rows, err := db.Pool.Query(ctx, query, pendingTimeInSeconds)
			if err != nil {
				return fmt.Errorf("worker ERR: selecting employees ERR: %w", err)
			}
			defer rows.Close()

			var checkedEmployees []*internal.Employee
			for rows.Next() {
				emp := internal.Employee{}

				err := rows.Scan(
					&emp.FirstName,
					&emp.LastName,
				)

				if err != nil {
					return fmt.Errorf("worker ERR: scanning employees ERR: %w", err)
				}

				checkedEmployees = append(checkedEmployees, &emp)
			}

			if len(checkedEmployees) > 0 {
				for _, emp := range checkedEmployees {
					if err := db.changeStatus(ctx, *emp); err != nil {
						return fmt.Errorf("worker ERR: changing status of %s %s ERR: %w",
							emp.FirstName, emp.LastName, err)
					}
				}
			}

		case <-ctx.Done():
			return nil
		}
	}
}
