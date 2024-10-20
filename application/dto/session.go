package dto

type CreateSessionOutputDto struct {
	QrCode  string `json:"qr_code"`
	Timeout int64  `json:"timeout"`
}

type CreateSessionInputDto struct {
	JID string `json:"jid"`
}

type DisconnectedSessionInputDto struct {
	JID string `json:"jid"`
}

type VerifySessionInputDto struct {
	JID string `json:"jid"`
}

type RetrieveSessionInputDto struct {
	JID string `json:"jid"`
}
