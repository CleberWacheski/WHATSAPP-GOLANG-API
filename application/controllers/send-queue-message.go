package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"whatsapp/application/dto"
	"whatsapp/application/queue"
	"whatsapp/application/utils"
)

func SendQueueMessage(w http.ResponseWriter, r *http.Request) {
	var body dto.SendQueueMessageInputDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.NewHttpError(w, err)
		return
	}
	r.Body.Close()
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
	if body.SecondsUntilProcessing <= 0 {
		utils.NewHttpError(w, errors.New("o tempo de delay não foi definido"))
		return
	}
	err = queue.QueueManager.NewTaskSendMessageQueue(queue.SendQueueMessageInput{
		JID:                    body.JID,
		Text:                   body.Text,
		SendJID:                body.SendJID,
		SecondsUntilProcessing: body.SecondsUntilProcessing,
	})
	if err != nil {
		utils.NewHttpError(w, err)
		return
	}
}
