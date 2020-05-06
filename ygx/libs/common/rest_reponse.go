package common

import "net/http"

// Response Status  `json:"data"`
type Response struct {
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type SearchListQuest struct {
	Page     int         `json:"page"`
	PageSize int         `json:"pagesize"`
	Data     interface{} `json:"data"`
}

type SearchListResponse struct {
	Count int         `json:"count"`
	Data  interface{} `json:"data"`
}

type YgxMsRequest func() (*http.Request, error)
type YgxMsReponse func() *Response

func (r *Response) GetErrorResponse(status int, message string) *Response {
	r.Status = status
	r.Message = message
	return r
}

type Message struct {
	StatusCode int32  `json:"StatusCode"`
	Message    string `json:"message"`
}

func (m *Message) GetMessage(messageCode string) {

}
