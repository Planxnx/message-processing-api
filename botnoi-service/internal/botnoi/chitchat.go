package botnoi

import (
	"encoding/json"
	"fmt"
	"strings"
)

type ChitChatResponse struct {
	Reply      string  `json:"reply"`
	Confidence float64 `json:"confidence"`
}

type chitchatService struct {
	Address string
	Token   string
}

// GetReplyMessage .
func (s *botnoiService) ChitChatMessage(rawMessage string) (string, error) {

	messageSlice := strings.Split(rawMessage, " ")
	message := strings.Join(messageSlice, "%20")

	endpoint := fmt.Sprintf("botnoichitchat?keyword=%s&styleid=1&botname=บอทใหญ่", message)
	resp, err := s.request(endpoint, nil)

	if err != nil {
		return "", err
	}

	responseBody := &ChitChatResponse{}
	json.Unmarshal(resp.Body(), responseBody)

	return responseBody.Reply, nil
}
