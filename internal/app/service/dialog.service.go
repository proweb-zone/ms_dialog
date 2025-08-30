package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ms_dialog/internal/app/dto"
	"ms_dialog/internal/app/entity"
	"ms_dialog/internal/app/repository"

	eventclient "github.com/proweb-zone/event-client"
	pb "github.com/proweb-zone/event-client/gen/go"
)

type DialogService struct {
	repo        *repository.DialogRepository
	eventClient *eventclient.EventClient
}

func NewDialogService(newEventClient *eventclient.EventClient, newRepo *repository.DialogRepository) *DialogService {
	return &DialogService{
		repo:        newRepo,
		eventClient: newEventClient,
	}
}

func (d *DialogService) SendMsgUser(requestDialog *dto.DialogRequestDto) (*entity.Dialog, error) {
	response, err := d.repo.SendMsgUser(&entity.Dialog{
		User_id_sender:    requestDialog.User_id_sender,
		User_id_recipient: requestDialog.User_id_recipient,
		Msg:               requestDialog.Msg,
	})

	if err != nil {
		return nil, err
	}

	d.eventClient.Publish(context.Background(), &pb.Event{
		Type:   "dialog.send",
		Source: "dialog-service",
		Payload: []byte(fmt.Sprintf(`{"User_id_sender": %d, "User_id_recipient": %d, "Id": %d, "Msg": "%s"}`,
			requestDialog.User_id_sender,
			requestDialog.User_id_recipient,
			response.ID,
			requestDialog.Msg)),
	})

	return response, nil
}

type EventDialogReponse struct {
	Id    int
	State bool
}

func (d *DialogService) CheckMsg(event *pb.Event) {
	payload := event.GetPayload()

	parsedResponse := &EventDialogReponse{}
	if err := json.Unmarshal(payload, &parsedResponse); err != nil {
		log.Fatalf("Failed to parse payload JSON: %v", err)
	}

	if parsedResponse.State == true {
		d.repo.ActiveMsg(parsedResponse.Id)
	} else {
		d.repo.DeleteMsg(parsedResponse.Id)
	}

}

func (d *DialogService) GetDialogList(userIdSender int, userIdRecepient int) (*[]entity.Dialog, error) {
	return d.repo.GetDialogList(userIdSender, userIdRecepient)
}
