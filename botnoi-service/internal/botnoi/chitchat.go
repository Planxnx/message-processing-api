package botnoi

import (
	"encoding/json"
	"fmt"
	"log"
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

func (*BotnoiService) getChitChatParams(rawMessage string) string {
	messageSlice := strings.Split(rawMessage, " ")
	return fmt.Sprintf("keyword=%s&styleid=1&botname=บอทใหญ่", strings.Join(messageSlice, "%20"))
}

// GetReplyMessage .
func (s *BotnoiService) ChitChatMessage(message string) (string, error) {
	endpoint := fmt.Sprintf("botnoichitchat?%s", s.getChitChatParams(message))
	resp, err := s.request(endpoint, nil)
	if err != nil {
		return "", err
	}

	responseBody := &ChitChatResponse{}
	json.Unmarshal(resp.Body(), responseBody)
	log.Println(responseBody)
	return responseBody.Reply, nil
}
