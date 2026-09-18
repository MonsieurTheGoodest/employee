package repository

import "time"

const (
	port                     = "5432"
	checkInterval            = 5 * time.Second
	pendingTimeInSeconds int = 180
	initSQLPath              = "/employee/internal/repository/initdb/schema.sql"
)
