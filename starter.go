package judger

import (
	"Rabbit-OJ-Backend/services/judger/config"
	"Rabbit-OJ-Backend/services/judger/mq"
	"Rabbit-OJ-Backend/services/judger/protobuf"
	"fmt"
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

type StarterType struct {
	Code       []byte
	IsContest  bool
	Sid        uint32
	Tid        uint32
	Version    uint32
	Language   string
	TimeLimit  uint32
	SpaceLimit uint32
	CompMode   string
}

func Starter(info *StarterType) error {
	request := &protobuf.JudgeRequest{
		Sid:        info.Sid,
		Tid:        info.Tid,
		Version:    strconv.FormatUint(uint64(info.Version), 10),
		Language:   info.Language,
		TimeLimit:  info.TimeLimit,
		SpaceLimit: info.SpaceLimit,
		CompMode:   info.CompMode,
		Code:       info.Code,
		Time:       time.Now().Unix(),
		IsContest:  info.IsContest,
	}

	pro, err := proto.Marshal(request)
	if err != nil {
		return err
	}

	if err := mq.PublishMessageSync(
		config.JudgeRequestTopicName,
		[]byte(fmt.Sprintf("%d%d", info.Sid, info.Tid)),
		pro); err != nil {
		return err
	}

	return nil
}
