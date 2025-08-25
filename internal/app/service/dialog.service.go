package service

import (
	"context"
	"fmt"
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
	fmt.Println("пользователь пытается отправить сообщение")

	response, err := d.repo.SendMsgUser(&entity.Dialog{
		User_id_sender:    requestDialog.User_id_sender,
		User_id_recipient: requestDialog.User_id_recipient,
		Msg:               requestDialog.Msg,
	})

	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}

	d.eventClient.Publish(context.Background(), &pb.Event{
		Type:   "dialog.send",
		Source: "dialog-service",
		Payload: []byte(fmt.Sprintf(`{
				"User_id_sender": %d,
        "User_id_recipient": %d,
        "msg": %s,
    }`, requestDialog.User_id_sender, requestDialog.User_id_recipient, requestDialog.Msg)),
	})

	return response, nil
}

// func (d *DialogService) GetDialogList(userIdSender int, userIdRecepient int) (*[]entity.Dialog, error) {
// 	return d.repo.GetDialogList(userIdSender, userIdRecepient)
// }
