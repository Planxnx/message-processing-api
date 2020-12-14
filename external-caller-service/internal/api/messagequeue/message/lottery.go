package message

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"

	"github.com/Planxnx/message-processing-api/external-caller-service/internal/lottery"
	messageschema "github.com/Planxnx/message-processing-api/message-schema"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const featureName string = "check-latest-lottery"

type CheckLatestLotteryRequestData struct {
	Number string `json:"number,omitempty"`
}

type CheckLatestLotteryAttachment struct {
	Date        string                    `json:"date,omitempty"`
	FoundReward []*CheckLatestLotteryData `json:"foundReward,omitempty"`
}

type CheckLatestLotteryData struct {
	Name   string `json:"name,omitempty"`
	Reward string `json:"reward,omitempty"`
	Number string `json:"number,omitempty"`
}

func (m *MessageHandler) CheckLatestLotteryHandler(msg *message.Message) error {
	defer msg.Ack()
	resultMsg := &messageschema.DefaultMessage{}
	proto.Unmarshal(msg.Payload, resultMsg)

	if strings.ToLower(resultMsg.Feature) != featureName {
		return nil
	}

	replymessage := &messageschema.DefaultMessage{
		Ref1:        resultMsg.Ref1,
		Ref2:        resultMsg.Ref2,
		Ref3:        resultMsg.Ref3,
		Owner:       resultMsg.Owner,
		PublishedBy: fmt.Sprintf("External Caller service: %v", featureName),
		Type:        "replyMessage",
	}

	//validate support mode
	if resultMsg.ExcuteMode != messageschema.ExecuteMode_Synchronous {
		log.Println("Wrong ExecMode")
		replymessage.Error = "wrong exec mode"
		replymessage.PublishedAt = timestamppb.Now()
		m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		return nil
	}

	requestData := &CheckLatestLotteryRequestData{}
	err := json.Unmarshal(resultMsg.Data, requestData)
	if err != nil {
		replymessage.ErrorInternal = err.Error()
		replymessage.PublishedAt = timestamppb.Now()
		m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		log.Printf("CheckLatestLotteryHandler Error: failed on get latest lotto: %v", err)
		return err
	}

	if requestData.Number == "" {
		replymessage.Error = "required lottery number"
		replymessage.PublishedAt = timestamppb.Now()
		m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		return nil
	}

	latestLotto, err := m.lotteryUsecase.GetLatestLottery()
	if err != nil {
		replymessage.ErrorInternal = err.Error()
		replymessage.PublishedAt = timestamppb.Now()
		m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)
		log.Printf("CheckLatestLotteryHandler Error: failed on get latest lotto: %v", err)
		return err
	}
	fmt.Println(requestData.Number)
	fmt.Println(latestLotto.Prizes[0].Number)

	foundReward := []*CheckLatestLotteryData{}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		var prizeWg sync.WaitGroup
		for _, prize := range latestLotto.Prizes {
			prizeWg.Add(1)
			go func(p lottery.LatestLotteryPrizes) {
				defer prizeWg.Done()
				i := sort.SearchStrings(p.Number, requestData.Number)
				if p.Number[i] == requestData.Number {
					foundReward = append(foundReward, &CheckLatestLotteryData{
						Name:   p.Name,
						Reward: p.Reward,
						Number: p.Number[i],
					})
				}
			}(prize)
		}
		prizeWg.Wait()
	}()
	go func() {
		defer wg.Done()
		var runningNumberWg sync.WaitGroup
		for _, runningNumbers := range latestLotto.RunningNumbers {
			runningNumberWg.Add(1)
			go func(p lottery.LatestLotteryRunningNumbers) {
				defer runningNumberWg.Done()
				i := sort.SearchStrings(p.Number, requestData.Number)
				if p.Number[i] == requestData.Number {
					foundReward = append(foundReward, &CheckLatestLotteryData{
						Name:   p.Name,
						Reward: p.Reward,
						Number: p.Number[i],
					})
				}
			}(runningNumbers)
		}
		runningNumberWg.Wait()
	}()
	wg.Wait()

	fmt.Println(foundReward)
	attachmentData, err := json.Marshal(&CheckLatestLotteryAttachment{
		Date:        latestLotto.Date,
		FoundReward: foundReward,
	})

	if err != nil {
		replymessage.ErrorInternal = err.Error()
		replymessage.PublishedAt = timestamppb.Now()
		m.messageUsecase.Emit(watermill.NewUUID(), resultMsg.CallbackTopic, replymessage)

		log.Printf("ChitchatHandler Error: failed on marshal attchment: %v", err)
		return err
	}

	replymessage.Data = attachmentData
	replymessage.PublishedAt = timestamppb.Now()
	replymessage.PublishedAt = timestamppb.Now()

	log.Printf("Replied !!%v\n", string(replymessage.Data))
	err = m.messageUsecase.Emit(watermill.NewShortUUID(), resultMsg.CallbackTopic, replymessage)
	if err != nil {
		log.Printf("ChitchatHandler Error: failed on emit message: %v", err)
		return err
	}
	return nil

}
