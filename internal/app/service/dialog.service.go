package service

import (
	"context"
	"fmt"

	eventclient "github.com/proweb-zone/event-client"
	pb "github.com/proweb-zone/event-client/gen/go"
)

type DialogService struct {
	//repo *repository.DialogRepository
	eventClient *eventclient.EventClient
}

func NewDialogService(newEventClient *eventclient.EventClient) *DialogService {
	return &DialogService{
		eventClient: newEventClient,
	}
}

func (d *DialogService) SendMsgUser() {
	fmt.Println("пользователь пытается отправить сообщение")
	d.eventClient.Publish(context.Background(), &pb.Event{
		Type:   "dialog.send",
		Source: "dialog-service",
		Payload: []byte(`{
        "user_id": "12345",
        "msg": "hello my friends",
    }`),
	})
	//	return d.repo.SendMsgUser(&entity.Dialog{
	//		User_id_sender:    requestDialog.User_id_sender,
	//		User_id_recipient: requestDialog.User_id_recipient,
	//		Msg:               requestDialog.Msg,
	//	})
}

// func (d *DialogService) GetDialogList(userIdSender int, userIdRecepient int) (*[]entity.Dialog, error) {
// 	return d.repo.GetDialogList(userIdSender, userIdRecepient)
// }
