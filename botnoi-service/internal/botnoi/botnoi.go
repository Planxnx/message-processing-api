package botnoi

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

type BotnoiService struct {
	address string
	token   string
}

func New(address string, token string) *BotnoiService {
	return &BotnoiService{
		address: address,
		token:   token,
	}
}

func (b *BotnoiService) request(endpoint string, body interface{}) (*fasthttp.Response, error) {
	url := fmt.Sprintf("%s/%s", b.address, endpoint)
	requestAuthorization := fmt.Sprintf("Bearer %s", b.token)
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetBody(requestBody)
	req.SetRequestURI(url)
	req.Header.Set("Authorization", requestAuthorization)
	req.Header.Set("Content-Type", "application/json")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, fmt.Errorf("botnoi request error: failed on send request: %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return resp, fmt.Errorf("botnoi request error:expected status code %d but  got %d", fasthttp.StatusOK, resp.StatusCode())
	}
	return resp, nil
}
