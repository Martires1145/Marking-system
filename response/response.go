package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var response = Response{
	Code:    200,
	Message: "",
	Data:    nil,
}

type PageResponse struct {
	PageSize int         `json:"pageSize"`
	Page     int         `json:"page"`
	HNext    bool        `json:"HNext"`
	Data     interface{} `json:"data"`
}

func Success(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	rsp := response
	rsp.Message = message
	rsp.Data = data
	r, _ := json.Marshal(rsp)
	_, _ = w.Write(r)
}

func Fail(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	rsp := response
	rsp.Code = code
	rsp.Message = message
	r, _ := json.Marshal(rsp)
	_, _ = w.Write(r)
}
