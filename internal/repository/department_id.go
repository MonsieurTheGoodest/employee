package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func createDepartment(
	ctx context.Context,
	department string,
	tx *pgx.Tx,
) (int, error) {

	query := `
		INSERT INTO departments (department)
		VALUES ($1)
		RETURNING id
	`

	var id int
	err := (*tx).QueryRow(ctx, query, department).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("scanning department ERR: %w", err)
	}

	return id, nil
}

func (db *DataBase) DepartmentID(
	ctx context.Context,
	department string,
) (int, error) {

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
	if err != nil {
		return 0, fmt.Errorf("starting transaction ERR: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		SELECT id
		FROM departments
		WHERE department = $1
	`

	rows, err := tx.Query(ctx, query, department)
	if err != nil {
		return 0, fmt.Errorf("receiving department ERR: %w", err)
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		if id != 0 {
			return 0, fmt.Errorf(`something went wrong:
				more than one department ERR`)
		}

		err := rows.Scan(&id)

		if err != nil {
			return 0, fmt.Errorf("scanning department ERR: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("reading department ERR: %w", err)
	}

	if id == 0 {
		id, err = createDepartment(ctx, department, &tx)

		if err != nil {
			return 0, fmt.Errorf("creating department ERR: %w", err)
		}
	}

	return id, tx.Commit(ctx)
}
