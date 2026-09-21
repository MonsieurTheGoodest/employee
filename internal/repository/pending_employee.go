package repository

import (
	"context"
	"employee/internal"
	"fmt"
)

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
