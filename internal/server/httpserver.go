package server

import (
	"context"
	"fmt"
	"log"
	"ms_dialog/internal/app/handlers"
	"ms_dialog/internal/app/repository"
	"ms_dialog/internal/app/routing"
	"ms_dialog/internal/app/service"
	"ms_dialog/internal/config"
	"ms_dialog/internal/db/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	eventclient "github.com/proweb-zone/event-client"
)

func StartServer(config *config.Config) {

	// connect event-client
	client, err := eventclient.New(eventclient.Config{
		GatewayAddress: config.GrpcServer.Addr,
		ServiceName:    "dialog-service",
		MaxRetries:     5,
		RetryDelay:     1 * time.Second,
	})

	if err != nil {
		log.Fatalf("Failed to create event client: %v", err)
	}

	defer client.Close()

	log.Println("MS Dialog service started")

	// init service dialog
	conn := postgres.Connect(config)
	defer postgres.Close(conn)

	dialogRepository := repository.NewDialogRepository(conn)
	newDialogService := service.NewDialogService(client, dialogRepository)

	authRepository := repository.NewAuthRepository(conn)
	newAuthService := service.NewAuthService(client, authRepository)

	// init handler
	handlers, err := handlers.NewHandlers(newDialogService, newAuthService)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	routing := routing.NewRouting(handlers)

	// subscribe on event
	err = client.Subscribe(context.Background(), []string{
		"user.access",
		"user.auth",
	}, handlers.Events)

	if err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	go func() {
		http.ListenAndServe(":"+config.HTTPServer.ServerPort, routing)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}
