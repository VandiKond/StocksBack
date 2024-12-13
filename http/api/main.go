package api

import (
	"encoding/json"
	"net/http"

	"github.com/vandi37/StocksBack/http/api/responses"
)

type Response struct {
	Ok          bool   `json:"ok"`
	StatusCode  int    `json:"status_code"`
	Description string `json:"description"`
	ContentType string `json:"content-type"`
	Data        any    `json:"data"`
}

func (r Response) Send(w http.ResponseWriter) error {
	w.WriteHeader(r.StatusCode)
	return json.NewEncoder(w).Encode(r)
}

func SendOkResponse(w http.ResponseWriter, data any, contentType string) error {
	resp := Response{
		Ok:          true,
		StatusCode:  http.StatusOK,
		Description: "OK",
		ContentType: contentType,
		Data:        data,
	}
	return resp.Send(w)
}

func SendErrorResponse(w http.ResponseWriter, status int, err error) error {
	resp := Response{
		Ok:          false,
		StatusCode:  status,
		Description: http.StatusText(status),
		ContentType: responses.ErrorType,
		Data:        err.Error(),
	}
	return resp.Send(w)
}
