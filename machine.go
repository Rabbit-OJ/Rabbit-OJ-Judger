package judger

import (
	"context"
	"github.com/Rabbit-OJ/Rabbit-OJ-Judger/logger"
	"sync"
)

var (
	MachineWaitGroup sync.WaitGroup
)

func StartMachine(ctx context.Context, index uint, queueChan chan []byte) {
	logger.Printf("[Machine] Concurrent #%d started \n", index)
	MachineWaitGroup.Add(1)
	defer MachineWaitGroup.Done()

	for {
		select {
		case delivery := <-queueChan:
			logger.Printf("[Machine] #%d machine START \n", index)
			JudgeRequestBridge(delivery)
			logger.Printf("[Machine] #%d machine FINISH \n", index)
		case <-ctx.Done():
			logger.Printf("[Machine] #%d machine Exited \n", index)
			return
		}
	}
}
