package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func (db *DataBase) CheckPendingEmployees(
	ctx context.Context,
	pendingTimeInSeconds int,
) ([]int, error) {
	query := `
		SELECT employee_id
		FROM pending_employees 
		WHERE creating_time < NOW() - $1 * INTERVAL '1 second'
	`

	rows, err := db.Pool.Query(ctx, query, pendingTimeInSeconds)
	if err != nil {
		return []int{}, fmt.Errorf("worker ERR: selecting employees ERR: %w", err)
	}
	defer rows.Close()

	var checkedEmployees []int
	for rows.Next() {
		var id int

		err := rows.Scan(&id)

		if err != nil {
			return []int{}, fmt.Errorf("worker ERR: scanning employees ERR: %w", err)
		}

		checkedEmployees = append(checkedEmployees, id)
	}

	return checkedEmployees, nil
}

func (db *DataBase) ChangeStatus(ctx context.Context, id int) error {
	checkedID, err := statusID(ctx, checked, db)
	if err != nil {
		return err
	}

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("starting transaction ERR: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE employees
		SET status_id = $1
		WHERE id = $2
	`

	_, err = tx.Exec(ctx, query, checkedID, id)

	if err != nil {
		return fmt.Errorf("updating status ERR: %w", err)
	}

	query = `
		DELETE FROM pending_employees
		WHERE employee_id = $1
	`

	_, err = tx.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("deleting employee ERR: %w", err)
	}

	return tx.Commit(ctx)
}
