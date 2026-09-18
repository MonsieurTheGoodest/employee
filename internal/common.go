package internal

import "fmt"

var (
	EmployeeAlreadyExists  error = fmt.Errorf("employee already exists")
	DepartmentDoesNotExist error = fmt.Errorf("department doesn't exist")
)

type Employee struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type EmployeeWithID struct {
	FullName string `json:"full_name"`
	ID       int    `json:"id"`
}

type EmployeesList struct {
	Department string `json:"department"`
	Count      int    `json:"count"`
	List       []*EmployeeWithID
}
