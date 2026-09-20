package repository

import (
	"context"
	"employee/internal"
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

func (fullEmp *fullEmployee) FullName() string {
	return fullEmp.LastName + " " + fullEmp.FirstName
}

func convertToEmployeesList(fullEmp []*fullEmployee) internal.EmployeesList {
	employeesList := internal.EmployeesList{}

	for i := 0; i < len(fullEmp); i++ {
		empWithID := &internal.EmployeeWithID{
			ID:       fullEmp[i].ID,
			FullName: fullEmp[i].FullName(),
		}

		employeesList.List = append(employeesList.List, empWithID)
	}

	employeesList.Count = len(employeesList.List)

	return employeesList
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

func departmentID(
	ctx context.Context,
	department string,
	tx *pgx.Tx,
) (int, error) {

	query := `
		SELECT id
		FROM departments
		WHERE department = $1
	`

	rows, err := (*tx).Query(ctx, query, department)
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
		id, err = createDepartment(ctx, department, tx)

		if err != nil {
			return 0, fmt.Errorf("creating department ERR: %w", err)
		}
	}

	return id, nil
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
			return 0, fmt.Errorf(`something went wrong:
				more than one same statuses ERR`)
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
		return true, fmt.Errorf(`something went wrong: 
				more than one employee with the same first and second names`)
	}

	return counter == 1, nil
}

func checkExistenceDepartment(
	ctx context.Context,
	department string,
	db *DataBase,
) (bool, error) {

	query := `
		SELECT id
		FROM departments
		WHERE department = $1
	`

	rows, err := db.Pool.Query(ctx, query, department)
	if err != nil {
		return false, fmt.Errorf("receiving department ERR: %w", err)
	}
	defer rows.Close()

	var counter int
	for rows.Next() {
		counter++
	}

	if err := rows.Err(); err != nil {
		return false, fmt.Errorf("reading department ERR: %w", err)
	}

	if counter >= 2 {
		return true, fmt.Errorf(`something went wrong: 
				more than one department`)
	}

	return counter == 1, nil
}

func (db *DataBase) CreateEmployee(
	ctx context.Context,
	emp internal.Employee,
	department string,
) error {

	pendingID, err := statusID(ctx, pending, db)
	if err != nil {
		return err
	}

	tx, err := db.Pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.Serializable,
	})
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

	depID, err := departmentID(ctx, department, &tx)
	if err != nil {
		return err
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

func (db *DataBase) EmployeesFromDepartment(
	ctx context.Context,
	department string,
) (internal.EmployeesList, error) {

	exists, err := checkExistenceDepartment(ctx, department, db)

	if err != nil {
		return internal.EmployeesList{}, err
	}

	if !exists {
		return internal.EmployeesList{}, internal.DepartmentDoesNotExist
	}

	query := `
		SELECT employees.first_name, employees.last_name, employees.id, departments.department, statuses.status
		FROM employees
		JOIN departments ON employees.department_id = departments.id
		JOIN statuses ON employees.status_id = statuses.id
		WHERE department = $1 AND status = $2
	`

	rows, err := db.Pool.Query(ctx, query, department, checked)
	if err != nil {
		return internal.EmployeesList{}, fmt.Errorf("receiving employees ERR: %w", err)
	}
	defer rows.Close()

	fullEmp := make([]*fullEmployee, 0)

	for rows.Next() {
		var emp fullEmployee

		err := rows.Scan(
			&emp.FirstName,
			&emp.LastName,
			&emp.ID,
			&emp.Department,
			&emp.Status,
		)

		if err != nil {
			return internal.EmployeesList{}, fmt.Errorf("scanning employess ERR: %w", err)
		}

		fullEmp = append(fullEmp, &emp)
	}

	err = rows.Err()

	if err != nil {
		return internal.EmployeesList{}, fmt.Errorf("reading employess ERR: %w", err)
	}

	return convertToEmployeesList(fullEmp), nil
}

func (db *DataBase) PendingEmployees(
	ctx context.Context,
) (internal.EmployeesList, error) {

	query := `
		SELECT employees.first_name, employees.last_name, employees.id, departments.department
		FROM pending_employees
		JOIN employees ON pending_employees.employee_id = employees.id
		JOIN departments ON employees.department_id = departments.id
	`

	rows, err := db.Pool.Query(ctx, query)
	if err != nil {
		return internal.EmployeesList{}, fmt.Errorf("receiving employees ERR: %w", err)
	}
	defer rows.Close()

	fullEmp := make([]*fullEmployee, 0)

	for rows.Next() {
		emp := fullEmployee{
			Status: pending,
		}

		err := rows.Scan(
			&emp.FirstName,
			&emp.LastName,
			&emp.ID,
			&emp.Department,
		)

		if err != nil {
			return internal.EmployeesList{}, fmt.Errorf("scanning employees ERR: %w", err)
		}

		fullEmp = append(fullEmp, &emp)
	}

	err = rows.Err()

	if err != nil {
		return internal.EmployeesList{}, fmt.Errorf("reading employess ERR: %w", err)
	}

	return convertToEmployeesList(fullEmp), nil
}

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
