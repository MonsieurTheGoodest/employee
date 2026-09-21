package repository

import (
	"context"
	"employee/internal"
	"fmt"
)

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
