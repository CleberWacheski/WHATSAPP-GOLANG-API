package controllers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"whatsapp/application/client"
	"whatsapp/application/utils"
)

func SendDocument(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.NewHttpError(w, err)
		return
	}
	if header.Filename == "" {
		utils.NewHttpError(w, errors.New("file is empty"))
		return
	}
	jid := r.Form.Get("jid")
	sendJID := r.Form.Get("send_jid")
	fileMsg := r.Form.Get("file_msg")
	if jid == "" {
		utils.NewHttpError(w, errors.New("file is empty"))
		return
	}
	if sendJID == "" {
		utils.NewHttpError(w, errors.New("file is empty"))
		return
	}
	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		if err != nil {
			utils.NewHttpError(w, errors.New("mime/type is not defined"))
			return
		}
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		utils.NewHttpError(w, err)
		return
	}
	fileBytes := buf.Bytes()
	client.WhatsappAPI.SendDocument(context.TODO(), jid, sendJID, fileBytes, header.Filename, mimeType, fileMsg)
	w.WriteHeader(200)
}
