package client

import (
	"context"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

type whatsmeowAPI struct {
}

var WhatsappAPI = whatsmeowAPI{}

var WhatsAppClient = make(map[string]*whatsmeow.Client)
var WhatsAppClientIsConnecting = make(map[string]bool)

var database *sqlstore.Container

type CreateSessionResponse struct {
	QrCode  string `json:"qr_code"`
	Timeout int64  `json:"timeout"`
}

type VerifyConnectedResponse struct {
	Connected bool `json:"connected"`
}

func (*whatsmeowAPI) Initialize() error {
	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", "file:prod.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}
	database = container
	devicesStores, err := container.GetAllDevices()
	if err != nil {
		return err
	}
	for _, v := range devicesStores {
		jid := fmt.Sprintf("%s@s.whatsapp.net", v.ID.User)
		client := whatsmeow.NewClient(v, nil)
		err = client.Connect()
		if err != nil {
			fmt.Println(err)
		}
		if err == nil {
			WhatsAppClient[jid] = client
		}
	}
	return nil
}

func eventHandler(evt interface{}, jid string) {
	switch evt.(type) {
	case *events.Connected:
		{
			WhatsAppClientIsConnecting[jid] = false
		}
	}
}

func (*whatsmeowAPI) CreateSession(ctx context.Context, jid string) (*CreateSessionResponse, error) {
	conn, isValid := WhatsAppClient[jid]
	if isValid {
		if conn.IsConnected() {
			conn.Disconnect()
		}
	}
	deviceStore := database.NewDevice()
	client := whatsmeow.NewClient(deviceStore, nil)
	qrChan, _ := client.GetQRChannel(ctx)
	client.AddEventHandler(func(evt interface{}) {
		eventHandler(evt, jid)
	})
	err := client.Connect()
	if err != nil {
		return nil, err
	}
	for evt := range qrChan {
		if evt.Event == "code" {
			WhatsAppClient[jid] = client
			WhatsAppClientIsConnecting[jid] = true
			return &CreateSessionResponse{
				QrCode:  evt.Code,
				Timeout: evt.Timeout.Milliseconds(),
			}, nil
		} else if evt.Event == "timeout" || evt.Event == "error" {
			return nil, fmt.Errorf("falha ao gerar o qr code, error: %s", evt.Error.Error())
		}
	}
	return nil, fmt.Errorf("erro inesperado ao gerar a sessão")
}

func (*whatsmeowAPI) DisconnectSession(ctx context.Context, jid string) error {
	conn, isValid := WhatsAppClient[jid]
	if isValid {
		if conn.IsConnected() {
			delete(WhatsAppClient, jid)
			conn.Disconnect()
		}
	}
	return nil
}

func (*whatsmeowAPI) VerifyConnected(ctx context.Context, jid string) *VerifyConnectedResponse {
	isConnecting := WhatsAppClientIsConnecting[jid]
	if isConnecting {
		return &VerifyConnectedResponse{
			Connected: false,
		}
	}
	conn, isValid := WhatsAppClient[jid]
	if !isValid {
		return &VerifyConnectedResponse{
			Connected: false,
		}
	}
	if !conn.IsConnected() {
		return &VerifyConnectedResponse{
			Connected: false,
		}
	}
	return &VerifyConnectedResponse{
		Connected: true,
	}
}

func (*whatsmeowAPI) SendMessage(ctx context.Context, jid string, toJid string, message string) error {
	isConnecting := WhatsAppClientIsConnecting[jid]
	if isConnecting {
		return errors.New("whatsapp não conectado")
	}
	conn, isValid := WhatsAppClient[jid]
	if !isValid || !conn.IsConnected() {
		return errors.New("whatsapp não conectado")
	}
	to, _ := types.ParseJID(toJid)
	log.Printf("enviando mensagem com jid=%s para o jid=%s a mensagem =%s", jid, toJid, message)
	_, err := conn.SendMessage(ctx, to, &waE2E.Message{
		Conversation: &message,
	})
	return err
}

func (*whatsmeowAPI) SendDocument(ctx context.Context, jid string, toJid string, file []byte, fileName string, mimetype string, fileMsg string) error {
	conn, isValid := WhatsAppClient[jid]
	if !isValid || !conn.IsConnected() {
		return errors.New("whatsapp não conectado")
	}
	to, _ := types.ParseJID(toJid)
	uploadedDoc, err := conn.Upload(ctx, file, whatsmeow.MediaDocument)
	if err != nil {
		return err
	}
	doc := &waE2E.DocumentMessage{
		URL:           &uploadedDoc.URL,
		Mimetype:      &mimetype,
		FileName:      &fileName,               // Nome do arquivo a ser exibido
		FileLength:    &uploadedDoc.FileLength, // Tamanho do arquivo
		MediaKey:      uploadedDoc.MediaKey,
		FileEncSHA256: uploadedDoc.FileEncSHA256,
		FileSHA256:    uploadedDoc.FileSHA256,
		DirectPath:    &uploadedDoc.DirectPath,
	}
	_, err = conn.SendMessage(ctx, to, &waE2E.Message{
		DocumentMessage: doc,
	})
	if fileMsg != "" {
		_, err := conn.SendMessage(ctx, to, &waE2E.Message{
			Conversation: &fileMsg,
		})
		return err
	}
	return err
}
