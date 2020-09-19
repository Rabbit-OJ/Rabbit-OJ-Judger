package judger

import (
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/logger"
	JudgerModels "github.com/Rabbit-OJ/Rabbit-OJ-Judger/models"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/mq"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/protobuf"
	"github.com/golang/protobuf/proto"
	"sync"
)

type JudgeResponseCallback = func(sid uint32, isContest bool, judgeResult []*JudgerModels.JudgeResult)

var (
	CallbackWaitGroup sync.WaitGroup

	OnJudgeResponse []JudgeResponseCallback
)

func CallbackAllError(status string, sid uint32, isContest bool, datasetCount int) {
	go func() {
		CallbackWaitGroup.Add(1)
		defer CallbackWaitGroup.Done()

		logger.Printf("(%d) Callback judge error with status: %s \n", sid, status)
		ceResult := make([]*protobuf.JudgeCaseResult, datasetCount)
		for i := range ceResult {
			ceResult[i] = &protobuf.JudgeCaseResult{
				Status: status,
			}
		}

		response := &protobuf.JudgeResponse{
			Sid:       sid,
			Result:    ceResult,
			IsContest: isContest,
		}

		pro, err := proto.Marshal(response)
		if err != nil {
			logger.Println(err)
			return
		}

		if err := mq.PublishMessageSync(config.JudgeResponseTopicName, []byte(fmt.Sprintf("%d", sid)), pro);
			err != nil {
			logger.Printf("[Callback] Error when callback error message to queue %+v \n", err)
		}
	}()
}

func CallbackSuccess(sid uint32, isContest bool, resultList []*protobuf.JudgeCaseResult) {
	go func() {
		CallbackWaitGroup.Add(1)
		defer CallbackWaitGroup.Done()

		logger.Printf("(%d) Callback judge success \n", sid)

		response := &protobuf.JudgeResponse{
			Sid:       sid,
			Result:    resultList,
			IsContest: isContest,
		}

		pro, err := proto.Marshal(response)
		if err != nil {
			logger.Println(err)
			return
		}

		if err := mq.PublishMessageSync(config.JudgeResponseTopicName, []byte(fmt.Sprintf("%d", sid)), pro);
			err != nil {
			logger.Printf("[Callback] Error when callback success message to queue %+v \n", err)
		}
	}()
}
