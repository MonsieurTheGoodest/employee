package repository

import (
	"context"
	"employee/internal"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func checkExistenceEmployee(
	ctx context.Context,
	emp internal.Employee,
	tx *pgx.Tx,
) (bool, error) {

	query := `
		SELECT first_name
		FROM employees
		WHERE first_name = $1 AND last_name = $2
	`

	rows, err := (*tx).Query(ctx, query, emp.FirstName, emp.LastName)
	if err != nil {
		return false, fmt.Errorf("receiving employee ERR: %w", err)
	}
	defer rows.Close()

	var counter int
	for rows.Next() {
		counter++
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("reading employee ERR: %w", err)
	}

	if counter >= 2 {
		return true, fmt.Errorf("something went wrong: " +
			"more than one employee with the same first and second names")
	}

	return counter == 1, nil
}

func (db *DataBase) CreateEmployee(
	ctx context.Context,
	emp internal.Employee,
	depID int,
) error {

	pendingID, err := statusID(ctx, pending, db)
	if err != nil {
		return err
	}

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("starting transaction ERR: %w", err)
	}
	defer tx.Rollback(ctx)

	exists, err := checkExistenceEmployee(ctx, emp, &tx)

	if err != nil {
		return err
	}

	if exists {
		return internal.EmployeeAlreadyExists
	}

	query := `
		INSERT INTO employees (first_name, last_name, department_id, status_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int

	err = tx.QueryRow(ctx,
		query,
		emp.FirstName,
		emp.LastName,
		depID,
		pendingID).Scan(&id)

	if err != nil {
		return fmt.Errorf("creating employee ERR: %w", err)
	}

	query = `
		INSERT INTO pending_employees (employee_id)
		VALUES ($1)
	`
	_, err = tx.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("creating pending_employee ERR: %w", err)
	}

	return tx.Commit(ctx)
}
