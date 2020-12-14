package lottery

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type LatestLotteryResponse struct {
	Status string            `json:"status"`
	Data   LatestLotteryData `json:"response"`
}
type LatestLotteryData struct {
	Date           string                        `json:"date"`
	Endpoint       string                        `json:"endpoint"`
	Prizes         []LatestLotteryPrizes         `json:"prizes"`
	RunningNumbers []LatestLotteryRunningNumbers `json:"runningNumbers"`
}
type LatestLotteryPrizes struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Reward string   `json:"reward"`
	Amount int      `json:"amount"`
	Number []string `json:"number"`
}
type LatestLotteryRunningNumbers struct {
	ID     string   `json:"id"`
	Name   string   `json:"name"`
	Reward string   `json:"reward"`
	Amount int      `json:"amount"`
	Number []string `json:"number"`
}

// GetLatestLottery .
func (s *LotteryUsecase) GetLatestLottery() (*LatestLotteryData, error) {
	endpoint := fmt.Sprintf("latest")
	resp, err := s.request(endpoint)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	responseBody := &LatestLotteryResponse{}
	if err := json.Unmarshal(resp.Body(), responseBody); err != nil {
		return nil, errors.WithStack(err)
	}

	return &responseBody.Data, nil
}
