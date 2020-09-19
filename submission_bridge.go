package judger

import (
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/logger"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/mq"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/protobuf"
	"time"

	"github.com/golang/protobuf/proto"
)

func JudgeRequestBridge(body []byte) {
	judgeRequest := &protobuf.JudgeRequest{}
	if err := proto.Unmarshal(body, judgeRequest); err != nil {
		logger.Println(err)
		return
	}

	if config.Global.Extensions.Expire.Enabled &&
		judgeRequest.Time-time.Now().Unix() > config.Global.Extensions.CheckJudge.Interval*int64(time.Minute) {
		logger.Printf("[Bridge] Received expired judge %d , will ignore this\n", judgeRequest.Sid)
		return
	}

	status, report, err := Scheduler(judgeRequest)
	sid := judgeRequest.Sid
	if status == "Internal Error" {
		logger.Printf("(%d) [Bridge] Requeued due to %+v \n", sid, err)
		Requeue(config.JudgeRequestTopicName, body)
	} else if status != "OK" {
		logger.Printf("(%d) [Bridge] Calling back results \n", judgeRequest.Sid)
		CallbackAllError(status, sid, judgeRequest.IsContest, len(report))
	} else if status == "OK" {
		logger.Printf("(%d) [Bridge] Calling back results \n", judgeRequest.Sid)
		CallbackSuccess(sid, judgeRequest.IsContest, report)
	}
}

func JudgeResponseBridge(body []byte) {
	judgeResult := &protobuf.JudgeResponse{}
	if err := proto.Unmarshal(body, judgeResult); err != nil {
		logger.Println(err)
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
