package response

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
	Status        string    `json:"status"`
}

type MessageBodyModel struct {
	Success bool        `json:"success"`
	Result  *string     `json:"result,omitempty"`
	Error   *ErrorModel `json:"error,omitempty"`
}

type ErrorModel struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
