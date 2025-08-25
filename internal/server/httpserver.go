package server

import (
	"context"
	"fmt"
	"log"
	"ms_dialog/internal/app/handlers"
	"ms_dialog/internal/app/repository"
	"ms_dialog/internal/app/service"
	"ms_dialog/internal/config"
	"ms_dialog/internal/db/postgres"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"

	eventclient "github.com/proweb-zone/event-client"
	pb "github.com/proweb-zone/event-client/gen/go"
)

func StartServer(config *config.Config) {

	// connect event client
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

	// subscribe on event
	err = client.Subscribe(context.Background(), []string{
		"getuser.info",
		"dialog.send",
	}, handleEvent)

	if err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	// init service dialog
	conn := postgres.Connect(config)
	defer postgres.Close(conn)

	dialogRepository := repository.InitDialogRepository(conn)
	newDialogService := service.NewDialogService(client, dialogRepository)

	// init handlers
	handlers, err := handlers.Init(newDialogService)
	if err != nil {
		fmt.Errorf("%v", err)
	}

	r := chi.NewRouter()

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Post("/v2/dialog/{user_id}/send", handlers.SendMsgUser)
	r.Get("/v2/dialog/{user_id}/list", handlers.GetDialog)

	go func() {
		http.ListenAndServe(":"+config.HTTPServer.ServerPort, r)
	}()

	// go func() {
	// 	for {
	// 		if err != nil {
	// 			log.Printf("Publish failed: %v", err)
	// 		}
	// 		time.Sleep(10 * time.Second)
	// 	}
	// }()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}

func handleEvent(event *pb.Event) error {
	fmt.Println(event)
	return nil
}
