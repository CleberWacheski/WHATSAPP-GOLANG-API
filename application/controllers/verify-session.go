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

func VerifySession(w http.ResponseWriter, r *http.Request) {
	var body dto.VerifySessionInputDto
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
	output := client.WhatsappAPI.VerifyConnected(context.TODO(), body.JID)
	utils.HttpJsonResponse(w, output)
}
