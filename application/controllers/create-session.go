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

func CreateSession(w http.ResponseWriter, r *http.Request) {
	var body dto.CreateSessionInputDto
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
	output, err := client.WhatsappAPI.CreateSession(context.TODO(), body.JID)
	if err != nil {
		utils.NewHttpError(w, err)
		return
	}
	utils.HttpJsonResponse(w, output)
}
