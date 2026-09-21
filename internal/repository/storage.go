package repository

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type fullEmployee struct {
	ID         int    `db:"id"`
	FirstName  string `db:"first_name"`
	LastName   string `db:"last_name"`
	Department string `db:"department"`
	Status     string `db:"status"`
}

type DataBase struct {
	Pool *pgxpool.Pool
}

func NewDatabase(ctx context.Context, databaseURL string) (*DataBase, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("creating config ERR: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)

	if err != nil {
		return nil, fmt.Errorf("creating pool ERR: %w", err)
	}

	db := &DataBase{Pool: pool}

	if err := db.initSchema(ctx); err != nil {
		db.Pool.Close()

		return nil, fmt.Errorf("init schema ERR: %w", err)
	}

	return db, nil
}

func (dataBase *DataBase) Close() {
	dataBase.Pool.Close()
}

func createStatus(ctx context.Context, status string, tx *pgx.Tx) error {
	query := `
		INSERT INTO statuses (status)
		VALUES ($1)
		ON CONFLICT DO NOTHING
	`
	_, err := (*tx).Exec(ctx, query, status)

	if err != nil {
		return fmt.Errorf("creating status ERR: %w", err)
	}

	return nil
}

func (db *DataBase) initSchema(ctx context.Context) error {
	content, err := os.ReadFile(os.Getenv("INITDB_PATH"))
	if err != nil {
		return fmt.Errorf("cannot read initdb.sql: %w", err)
	}

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("cannot begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, string(content))
	if err != nil {
		return fmt.Errorf("schema execution ERR: %w", err)
	}

	for i := 0; i < len(statuses); i++ {
		createStatus(ctx, statuses[i], &tx)
	}

	return tx.Commit(ctx)
}

func statusID(ctx context.Context, status string, db *DataBase) (int, error) {
	query := `
		SELECT id
		FROM statuses
		WHERE status = $1
	`

	rows, err := db.Pool.Query(ctx, query, status)
	if err != nil {
		return 0, fmt.Errorf("receiving status ERR: %w", err)
	}
	defer rows.Close()

	var id int
	for rows.Next() {
		if id != 0 {
			return 0, fmt.Errorf("something went wrong: " +
				"more than one same statuses ERR")
		}

		err := rows.Scan(&id)

		if err != nil {
			return 0, fmt.Errorf("scanning status ERR: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return 0, fmt.Errorf("reading status ERR: %w", err)
	}

	return id, nil
}
