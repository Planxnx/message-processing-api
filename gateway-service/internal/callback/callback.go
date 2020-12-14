package callback

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type CallbackUsecase struct {
	address string
	token   string
}

func New() *CallbackUsecase {
	return &CallbackUsecase{}
}

func (*CallbackUsecase) Request(endpoint string, body interface{}) (*fasthttp.Response, error) {
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody(requestBody)
	req.SetRequestURI(endpoint)
	req.Header.Set("Content-Type", "application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("callback request error: failed on send request: %v", err)
	}
	return resp, nil
}
