package handlers

import (
	"fmt"
	"io"
	"log"
	"ms_dialog/internal/app/dto"
	"ms_dialog/internal/app/entity"
	"ms_dialog/internal/app/service"
	"ms_dialog/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi"
	pb "github.com/proweb-zone/event-client/gen/go"
)

type Handler struct {
	dialogService *service.DialogService
	authService   *service.AuthService
}

func NewHandlers(newDialogService *service.DialogService, newAthService *service.AuthService) (*Handler, error) {
	return &Handler{dialogService: newDialogService, authService: newAthService}, nil
}

func (h *Handler) Events(event *pb.Event) error {
	switch event.Type {
	case "user.access":
		h.dialogService.CheckMsg(event)
		return nil
	case "user.auth":
		h.authService.CreateToken(event)
		return nil
	default:
		log.Printf("Unknown event type: %s", event.Type)
		return nil
	}
}

func (h *Handler) SendMsgUser(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	// добавить
	// userId := 1

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

	msg, errSendMsgUser := h.dialogService.SendMsgUser(&requestDialogDto)
	if errSendMsgUser != nil {
		http.Error(w, errSendMsgUser.Error(), http.StatusBadRequest)
		return
	}

	utils.ResponseJson(msg, w)
	//w.Write([]byte("send msg user"))
	return
}

func (h *Handler) GetDialog(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userIdSender := auth.User_id
	// userIdSender := 1

	userIdRecepientStr := chi.URLParam(r, "user_id")
	userIdRecepient, err := strconv.Atoi(userIdRecepientStr)
	if err != nil {
		http.Error(w, "Error: User id "+userIdRecepientStr+"  не найден", http.StatusBadRequest)
		return
	}

	if userIdSender == userIdRecepient {
		http.Error(w, "Вы не можете получать диалог самого себя", http.StatusBadRequest)
		return
	}

	dialogList, errorDialog := h.dialogService.GetDialogList(userIdSender, userIdRecepient)
	if errorDialog != nil {
		http.Error(w, errorDialog.Error(), http.StatusBadRequest)
		return
	}

	utils.ResponseJson(dialogList, w)

	//w.Write([]byte("get dialog"))
	return
}

func (h *Handler) checkTokenAccess(r *http.Request) (*entity.Auth, error) {
	// Извлечение токена из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("Authorization header missing")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return nil, fmt.Errorf("Invalid authorization header format")
	}

	token := bearerToken[1]

	auth, err := h.authService.CheckAccessToken(token)
	if err != nil {
		return nil, err
	}

	return auth, nil
}
