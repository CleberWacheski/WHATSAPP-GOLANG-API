package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"whatsapp/application/client"
	"whatsapp/application/dto"
	"whatsapp/application/utils"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
	var body dto.SendMessageInputDto
	err := json.NewDecoder(r.Body).Decode(&body)
	r.Body.Close()
	if err != nil {
		utils.NewHttpError(w, err)
		return
	}
	if body.JID == "" {
		utils.NewHttpError(w, errors.New("jid não foi definido"))
		return
	}
	if body.SendJID == "" {
		utils.NewHttpError(w, errors.New("jid do contato não foi definido"))
		return
	}
	if body.Text == "" {
		utils.NewHttpError(w, errors.New("a mensagem não foi definida"))
		return
	}
	err = client.WhatsappAPI.SendMessage(context.TODO(), body.JID, body.SendJID, body.Text)
	if err != nil {
		utils.NewHttpError(w, err)
		return
	}
}
