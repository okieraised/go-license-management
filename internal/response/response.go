package response

import (
	"context"
	"fmt"
	"go-license-management/internal/constants"
	"time"
)

type Response struct {
	RequestID    string      `json:"request_id"`
	ErrorCode    string      `json:"code"`
	ErrorMessage string      `json:"message"`
	ServerTime   int64       `json:"server_time"`
	Count        int         `json:"count,omitempty"`
	Data         interface{} `json:"data"`
	Agg          interface{} `json:"agg,omitempty"`
	Meta         interface{} `json:"meta,omitempty"`
}

func NewResponse(ctx context.Context) *Response {
	resp := new(Response)
	resp.RequestID = fmt.Sprintf("%v", ctx.Value(constants.RequestIDField))
	resp.ServerTime = time.Now().Unix()
	resp.ErrorCode = "OK"
	resp.ErrorMessage = "OK"
	resp.Data = map[string]interface{}{}
	return resp
}

func (resp *Response) ToResponse(code, message string, data, meta, count interface{}) *Response {
	resp.ErrorCode = code
	resp.ErrorMessage = message

	if data != nil {
		resp.Data = data
	}

	if meta != nil {
		resp.Meta = meta
	}
	if count != nil {
		if _, ok := count.(int); ok {
			resp.Count = count.(int)
		}
	}
	return resp
}
