package repository

import "employee/internal"

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
