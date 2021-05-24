package request

import "time"

type MessageModel struct {
	Header MessageHeaderModel `json:"header"`
	Body   MessageBodyModel   `json:"body"`
}

type MessageHeaderModel struct {
	Title         string    `json:"title"`
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	Service       string    `json:"service"`
	CorrelationID string    `json:"correlationId"`
}

type MessageBodyModel struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
}
