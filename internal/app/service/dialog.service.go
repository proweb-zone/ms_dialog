package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"ms_dialog/internal/app/dto"
	"ms_dialog/internal/app/entity"
	"ms_dialog/internal/app/repository"
	"ms_dialog/internal/db/redis"

	eventclient "github.com/proweb-zone/event-client"
	pb "github.com/proweb-zone/event-client/gen/go"
)

type DialogService struct {
	repo        *repository.DialogRepository
	eventClient *eventclient.EventClient
	connRedisDb *redis.RedisDb
}

func NewDialogService(newEventClient *eventclient.EventClient, newRepo *repository.DialogRepository, newConnRedisDb *redis.RedisDb) *DialogService {
	return &DialogService{
		repo:        newRepo,
		eventClient: newEventClient,
		connRedisDb: newConnRedisDb,
	}
}

func (d *DialogService) SendMsgUser(requestDialog *dto.DialogRequestDto) (*entity.Dialog, error) {
	response, err := d.repo.SendMsgUser(&entity.Dialog{
		User_id_sender:    requestDialog.User_id_sender,
		User_id_recipient: requestDialog.User_id_recipient,
		Msg:               requestDialog.Msg,
		State:             false,
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
		dialog, errActiveMsg := d.repo.ActiveMsg(parsedResponse.Id)
		if errActiveMsg != nil {
			log.Fatalf("Failed to parse payload JSON: %v", errActiveMsg)
		}

		// add msg in redis
		errAddMsg := d.connRedisDb.AddMsg(dialog)
		if errAddMsg != nil {
			fmt.Errorf("error write msm in redis db")
		}

	} else {
		d.repo.DeleteMsg(parsedResponse.Id)
	}

}

func (d *DialogService) GetDialogList(userId int, userIdFriend int) (*[]entity.Dialog, error) {

	return d.connRedisDb.UdfGetMessages(userId, userIdFriend)

	//return d.repo.GetDialogList(userIdSender, userIdRecepient)
}
