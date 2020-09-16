package judger

import (
	"context"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/config"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/mq"
)

var (
	MachineContext           context.Context
	MachineContextCancelFunc context.CancelFunc
)

func JudgeRequestHandler() {
	queueChan := make(chan []byte)

	MachineContext, MachineContextCancelFunc = context.WithCancel(context.Background())
	for i := uint(0); i < config.Global.Concurrent.Judge; i++ {
		go StartMachine(MachineContext, i, queueChan)
	}

	for {
		select {
		case delivery := <-mq.JudgeRequestDeliveryChan:
			queueChan <- delivery
		case <-MachineContext.Done():
			return
		}
	}
}
