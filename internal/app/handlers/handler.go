package handlers

import (
	"io"
	"ms_dialog/internal/app/dto"
	"ms_dialog/internal/app/service"
	"ms_dialog/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type Handler struct {
	dialogService *service.DialogService
}

func Init(newDialogService *service.DialogService) (*Handler, error) {
	return &Handler{dialogService: newDialogService}, nil
}

func (h *Handler) SendMsgUser(w http.ResponseWriter, r *http.Request) {
	//h.dialogService.SendMsgUser()

	// auth, errAccessToken := h.checkTokenAccess(r)

	// if errAccessToken != nil {
	// 	http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
	// 	return
	// }

	// userId := auth.User_id

	userId := 2

	userIdRecepientStr := chi.URLParam(r, "user_id")
	userIdRecepient, err := strconv.Atoi(userIdRecepientStr)
	if err != nil {
		http.Error(w, "Error: User id "+userIdRecepientStr+"  не найден", http.StatusBadRequest)
		return
	}

	if userId == userIdRecepient {
		http.Error(w, "Вы не можете отправлять письмо самим себе", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestDialogDto dto.DialogRequestDto
	if err := utils.DecodeJson(body, &requestDialogDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestDialogDto.User_id_sender = userId
	requestDialogDto.User_id_recipient = userIdRecepient

	dialogId, errSendMsgUser := h.dialogService.SendMsgUser(&requestDialogDto)
	if errSendMsgUser != nil {
		http.Error(w, errSendMsgUser.Error(), http.StatusBadRequest)
		return
	}

	utils.ResponseJson(dialogId, w)
	//w.Write([]byte("send msg user"))
	return
}

func (h *Handler) GetDialog(w http.ResponseWriter, r *http.Request) {
	// auth, errAccessToken := h.checkTokenAccess(r)

	// if errAccessToken != nil {
	// 	http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
	// 	return
	// }

	// userIdSender := auth.User_id
	// userIdSender := 1

	// userIdRecepientStr := chi.URLParam(r, "user_id")
	// userIdRecepient, err := strconv.Atoi(userIdRecepientStr)
	// if err != nil {
	// 	http.Error(w, "Error: User id "+userIdRecepientStr+"  не найден", http.StatusBadRequest)
	// 	return
	// }

	// if userIdSender == userIdRecepient {
	// 	http.Error(w, "Вы не можете получать диалог самого себя", http.StatusBadRequest)
	// 	return
	// }

	// dialogList, errorDialog := h.dialogService.GetDialogList(userIdSender, userIdRecepient)
	// if errorDialog != nil {
	// 	http.Error(w, errorDialog.Error(), http.StatusBadRequest)
	// 	return
	// }

	// utils.ResponseJson(dialogList, w)

	w.Write([]byte("get dialog"))
	return
}
