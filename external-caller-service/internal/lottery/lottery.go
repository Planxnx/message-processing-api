package lottery

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type LotteryUsecase struct {
	address string
	token   string
}

func New(address string) *LotteryUsecase {
	return &LotteryUsecase{
		address: address,
	}
}

//request need to call fasthttp.ReleaseResponse() when finised use resp
func (b *LotteryUsecase) request(endpoint string) (*fasthttp.Response, error) {
	url := fmt.Sprintf("%s/%s", b.address, endpoint)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(url)
	req.Header.Set("Content-Type", "application/json")

	resp := fasthttp.AcquireResponse()

	if err := fasthttp.Do(req, resp); err != nil {
		return nil, errors.Errorf("lottery request error: failed on send request: %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return resp, errors.Errorf("lottery request error:expected status code %d but  got %d", fasthttp.StatusOK, resp.StatusCode())
	}
	return resp, nil
}
