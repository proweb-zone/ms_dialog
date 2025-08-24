package eventclient

func NewEventClient() {
	// Инициализация клиента
	// client, err := eventclient.New(eventclient.Config{
	// 	GatewayAddress: "localhost:50051",
	// 	ServiceName:    "dialog-service",
	// 	MaxRetries:     5,
	// 	RetryDelay:     1 * time.Second,
	// })

	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to create event client: %v", err)
	// }

	// defer client.Close()

	// log.Println("Dialog service started")

	// Ожидание сигнала завершения
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// <-sigChan
	// log.Println("Shutting down...")

}
