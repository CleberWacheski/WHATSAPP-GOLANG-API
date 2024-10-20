package dto

type SendMessageInputDto struct {
	JID     string `json:"jid"`
	Text    string `json:"text"`
	SendJID string `json:"send_jid"`
}

type SendQueueMessageInputDto struct {
	JID                    string `json:"jid"`
	Text                   string `json:"text"`
	SendJID                string `json:"send_jid"`
	SecondsUntilProcessing int64  `json:"seconds_until_processing"`
}
