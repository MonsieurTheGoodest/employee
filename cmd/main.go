package main

import (
	"context"
	"employee/internal/handler"
	"employee/internal/repository"
	"employee/internal/service"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := repository.NewDatabase(context.Background())

	if err != nil {
		panic(err)
	}

	defer db.Close()

	go func() {
		err := repository.Checking(db)

		if err != nil {
			log.Print(err)
		}
	}()

	service := service.NewService(db)
	handler := handler.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:         ":9000",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
