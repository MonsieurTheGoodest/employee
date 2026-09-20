package service

import (
	"context"
	"employee/internal"
	"employee/internal/repository"
	"fmt"
)

type Service struct {
	db *repository.DataBase
}

func NewService(db *repository.DataBase) *Service {
	return &Service{
		db: db,
	}
}

func (s *Service) CreateEmployee(
	ctx context.Context, emp internal.Employee, department string) error {

	depID, err := s.db.DepartmentID(ctx, department)

	if err != nil {
		return fmt.Errorf("create employee ERR:%s", err)
	}

	return s.db.CreateEmployee(ctx, emp, depID)
}

func (s *Service) EmployeesFromDepartment(
	ctx context.Context, department string) (internal.EmployeesList, error) {

	return s.db.EmployeesFromDepartment(ctx, department)
}

func (s *Service) PendingEmployees(
	ctx context.Context) (internal.EmployeesList, error) {

	return s.db.PendingEmployees(ctx)
}
