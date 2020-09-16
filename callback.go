package judger

import (
	"Rabbit-OJ-Backend/models"
	"Rabbit-OJ-Backend/services/judger/config"
	"Rabbit-OJ-Backend/services/judger/mq"
	"Rabbit-OJ-Backend/services/judger/protobuf"
	StorageService "Rabbit-OJ-Backend/services/storage"
	"fmt"
	"github.com/golang/protobuf/proto"
	"sync"
)

type JudgeResponseCallback = func(sid uint32, isContest bool, judgeResult []*models.JudgeResult)

var (
	CallbackWaitGroup sync.WaitGroup

	OnJudgeResponse []JudgeResponseCallback
)

func CallbackAllError(status string, sid uint32, isContest bool, storage *StorageService.Storage) {
	go func() {
		CallbackWaitGroup.Add(1)
		defer CallbackWaitGroup.Done()

		fmt.Printf("(%d) Callback judge error with status: %s \n", sid, status)

		ceResult := make([]*protobuf.JudgeCaseResult, storage.DatasetCount)
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

		mq.PublishMessage(config.JudgeResponseTopicName, []byte(fmt.Sprintf("%d", sid)), pro, true)
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

		mq.PublishMessage(config.JudgeResponseTopicName, []byte(fmt.Sprintf("%d", sid)), pro, true)
	}()
}
