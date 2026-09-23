package main

import (
	"context"
	"employee/config"
	"employee/internal/handler"
	"employee/internal/repository"
	"employee/internal/service"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
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

	service := service.NewService(db)

	serviceCtx, serviceCancel := context.WithCancel(context.Background())
	defer serviceCancel()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := service.Checking(serviceCtx, cfg.CheckInterval, cfg.PendingTimeInSeconds)

		if err != nil && !errors.Is(err, context.Canceled) {
			log.Print(err)
		}
	}()

	handler := handler.NewHandler(service)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := cfg.HTTPServer.Server()
	server.Handler = mux

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		err = server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server ERR: %s", err)
		}
	}()

	<-shutdownCtx.Done()

	stop()

	serviceCancel()

	serverCtx, serverCancel := context.WithTimeout(
		context.Background(),
		cfg.CloseServerTime,
	)
	defer serverCancel()

	err = server.Shutdown(serverCtx)

	if err != nil {
		log.Print(err)
	}

	serviceCtx, serviceCancel = context.WithTimeout(
		context.Background(),
		cfg.CloseServiceTime,
	)
	defer serviceCancel()

	workerDone := make(chan struct{})

	go func() {
		wg.Wait()
		close(workerDone)
	}()

	select {
	case <-workerDone:
		log.Println("worker stopped gracefully")

	case <-serviceCtx.Done():
		log.Println("worker shutdown timeout")
	}
}
