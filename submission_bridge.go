package judger

import (
	"Rabbit-OJ-Backend/services/judger/config"
	JudgerModels "Rabbit-OJ-Backend/services/judger/models"
	"Rabbit-OJ-Backend/services/judger/mq"
	"Rabbit-OJ-Backend/services/judger/protobuf"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
)

func JudgeRequestBridge(body []byte) {
	judgeRequest := &protobuf.JudgeRequest{}
	if err := proto.Unmarshal(body, judgeRequest); err != nil {
		fmt.Println(err)
		return
	}

	if config.Global.Extensions.Expire.Enabled &&
		judgeRequest.Time-time.Now().Unix() > config.Global.Extensions.CheckJudge.Interval*int64(time.Minute) {
		fmt.Printf("[Bridge] Received expired judge %d , will ignore this\n", judgeRequest.Sid)
		return
	}

	if alreadyAcked, err := Scheduler(judgeRequest); err != nil {
		if !alreadyAcked {
			Requeue(config.JudgeRequestTopicName, body)
		}

		fmt.Println(err)
		return
	}
}

func JudgeResponseBridge(body []byte) {
	judgeResult := &protobuf.JudgeResponse{}
	if err := proto.Unmarshal(body, judgeResult); err != nil {
		fmt.Println(err)
		return
	}

	judgeCaseResult := make([]*JudgerModels.JudgeResult, len(judgeResult.Result))
	for i, item := range judgeResult.Result {
		judgeCaseResult[i] = &JudgerModels.JudgeResult{
			Status:    item.Status,
			TimeUsed:  item.TimeUsed,
			SpaceUsed: item.SpaceUsed,
		}
	}

	for _, cb := range OnJudgeResponse {
		cb(judgeResult.Sid, judgeResult.IsContest, judgeCaseResult)
	}
}

func Requeue(topic string, body []byte) {
	mq.PublishMessageAsync(topic, []byte(fmt.Sprintf("%d", time.Now().UnixNano())), body)
}

func JudgeResultHandler() {
	for delivery := range mq.JudgeResponseDeliveryChan {
		go JudgeResponseBridge(delivery)
	}
}
