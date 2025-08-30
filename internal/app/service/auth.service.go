package service

import (
	"encoding/json"
	"log"
	"ms_dialog/internal/app/repository"

	eventclient "github.com/proweb-zone/event-client"
	pb "github.com/proweb-zone/event-client/gen/go"
)

type AuthService struct {
	repo        *repository.AuthRepository
	eventClient *eventclient.EventClient
}

type EventAuthReponse struct {
	User_id int
	Token   string
}

func NewAuthService(newEventClient *eventclient.EventClient, newRepo *repository.AuthRepository) *AuthService {
	return &AuthService{
		repo:        newRepo,
		eventClient: newEventClient,
	}
}

func (a *AuthService) CreateToken(event *pb.Event) error {
	payload := event.GetPayload()

	parsedResponse := &EventAuthReponse{}
	if err := json.Unmarshal(payload, &parsedResponse); err != nil {
		log.Fatalf("Failed to parse payload JSON: %v", err)
	}

	userId := parsedResponse.User_id
	token := parsedResponse.Token

	a.repo.CreateAuth(userId, token)

	return nil
}
