package judger

import (
	"fmt"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
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

func CallbackAllError(status string, sid uint32, isContest bool, datasetCount uint32) {
	go func() {
		CallbackWaitGroup.Add(1)
		defer CallbackWaitGroup.Done()

		fmt.Printf("(%d) Callback judge error with status: %s \n", sid, status)
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
			fmt.Println(err)
			return
		}

		mq.PublishMessageAsync(config.JudgeResponseTopicName, []byte(fmt.Sprintf("%d", sid)), pro)
	}()
}

func CallbackSuccess(sid uint32, isContest bool, resultList []*protobuf.JudgeCaseResult) {
	go func() {
		CallbackWaitGroup.Add(1)
		defer CallbackWaitGroup.Done()

		fmt.Printf("(%d) Callback judge success \n", sid)

		response := &protobuf.JudgeResponse{
			Sid:       sid,
			Result:    resultList,
			IsContest: isContest,
		}

		pro, err := proto.Marshal(response)
		if err != nil {
			fmt.Println(err)
			return
		}

		mq.PublishMessageAsync(config.JudgeResponseTopicName, []byte(fmt.Sprintf("%d", sid)), pro)
	}()
}
