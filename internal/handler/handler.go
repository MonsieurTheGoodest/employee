package handler

import (
	"employee/internal"
	"employee/internal/service"
	"encoding/json"
	"errors"
	"net/http"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /employee", h.createEmployee)
	mux.HandleFunc("GET /department/{department}", h.employeesFromDepartment)
	mux.HandleFunc("GET /pending", h.pendingEmployees)
}

func (h *Handler) createEmployee(w http.ResponseWriter, r *http.Request) {
	type inputEmployee struct {
		FirstName  string `json:"first_name"`
		LastName   string `json:"last_name"`
		Department string `json:"department"`
	}

	var inputEmp inputEmployee

	if err := json.NewDecoder(r.Body).Decode(&inputEmp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	emp := internal.Employee{
		FirstName: inputEmp.FirstName,
		LastName:  inputEmp.LastName,
	}

	err := h.service.CreateEmployee(r.Context(), emp, inputEmp.Department)

	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		if errors.Is(err, internal.EmployeeAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) employeesFromDepartment(w http.ResponseWriter, r *http.Request) {
	department := r.PathValue("department")

	employees, err := h.service.EmployeesFromDepartment(r.Context(), department)

	if err != nil {
		if errors.Is(err, internal.DepartmentDoesNotExist) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) pendingEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.service.PendingEmployees(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(employees); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
