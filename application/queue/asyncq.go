package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
	"whatsapp/application/client"
	"whatsapp/application/utils"

	"github.com/hibiken/asynq"
)

type queueManager struct {
}

var QueueManager = queueManager{}
var asyncqQueueClient *asynq.Client

type SendQueueMessageInput struct {
	JID                    string `json:"jid"`
	Text                   string `json:"text"`
	SendJID                string `json:"send_jid"`
	SecondsUntilProcessing int64  `json:"seconds_until_processing"`
}

func NewAsyncQueueManagerInitialize() {
	asyncqQueueClient = asynq.NewClient(asynq.RedisClientOpt{
		Addr: utils.ENV.REDIS_URL,
	})
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: utils.ENV.REDIS_URL},
		asynq.Config{
			Concurrency: 5,
		},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc("queue:message", sendQueueMessageProcess)
	if err := server.Run(mux); err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %v", err)
	}
}

func (q *queueManager) NewTaskSendMessageQueue(dto SendQueueMessageInput) error {
	jsonPayload, err := json.Marshal(dto)
	if err != nil {
		return err
	}
	if asyncqQueueClient == nil {
		return errors.New("queue client is not connected")
	}
	task := asynq.NewTask("queue:message", jsonPayload)
	p := time.Now().Add(time.Duration(dto.SecondsUntilProcessing * int64(time.Second)))
	_, err = asyncqQueueClient.Enqueue(task, asynq.ProcessAt(p))
	return err
}

func sendQueueMessageProcess(ctx context.Context, t *asynq.Task) error {
	var p SendQueueMessageInput
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("falha ao obter o json na fila de envio de mensagens")
	}
	log.Printf("Enviando mensagem(FILA) com jid=%s para o jid=%s a mensagem =%s", p.JID, p.SendJID, p.Text)
	client.WhatsappAPI.SendMessage(context.TODO(), p.JID, p.SendJID, p.Text)
	return nil
}
