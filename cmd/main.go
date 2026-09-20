package main

import (
	"context"
	"employee/config"
	"employee/internal/handler"
	"employee/internal/repository"
	"employee/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		panic(err)
	}

	db, err := repository.NewDatabase(context.Background(), cfg.StorageURL())

	if err != nil {
		panic(err)
	}

	defer db.Close()

	go func() {
		err := repository.Checking(db, cfg.CheckInterval, cfg.PendingTime)

		if err != nil {
			log.Print(err)
		}
	}()

	service := service.NewService(db)
	handler := handler.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := cfg.HTTPServer.Server()
	server.Handler = mux

	log.Fatal(server.ListenAndServe())
}
