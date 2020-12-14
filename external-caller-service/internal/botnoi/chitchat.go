package botnoi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type ChitChatResponse struct {
	Reply      string  `json:"reply"`
	Confidence float64 `json:"confidence"`
}

type chitchatService struct {
	Address string
	Token   string
}

func (*BotnoiService) getChitChatParams(rawMessage string) string {
	messageSlice := strings.Split(rawMessage, " ")
	return fmt.Sprintf("keyword=%s&styleid=1&botname=บอทใหญ่", strings.Join(messageSlice, "%20"))
}

// GetReplyMessage .
func (s *BotnoiService) ChitChatMessage(message string) (string, error) {
	endpoint := fmt.Sprintf("botnoichitchat?%s", s.getChitChatParams(message))
	resp, err := s.request(endpoint, nil)
	defer fasthttp.ReleaseResponse(resp)
	if err != nil {
		return "", errors.WithStack(err)
	}

	responseBody := &ChitChatResponse{}
	if err := json.Unmarshal(resp.Body(), responseBody); err != nil {
		return "", errors.WithStack(err)
	}

	return responseBody.Reply, nil
}
