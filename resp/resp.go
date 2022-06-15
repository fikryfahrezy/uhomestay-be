package resp

import (
	"encoding/json"
	"net/http"
	"strings"
)

type Response struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Error      error  `json:"-"`
}

func NewResponse(statusCode int, message string, err error) Response {
	if message == "" && err != nil {
		message = err.Error()
	}

	return Response{
		StatusCode: statusCode,
		Message:    message,
		Error:      err,
	}
}

func (r Response) HttpJSON(w http.ResponseWriter, body interface{}) {
	if ct := w.Header().Get("Content-Type"); ct != "" && !strings.Contains(ct, "application/json") {
		w.Header().Add("Content-Type", "application/json")
	}

	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", "application/json")
	}

	encoder := json.NewEncoder(w)
	w.WriteHeader(r.StatusCode)

	if r.Error == nil && body != nil {
		encoder.Encode(body)
		return
	}

	encoder.Encode(r)
}

type HttpBody struct {
	Data interface{} `json:"data"`
}

func NewHttpBody(data interface{}) HttpBody {
	return HttpBody{
		Data: data,
	}
}
