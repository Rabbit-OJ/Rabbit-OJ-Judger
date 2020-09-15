package judger

import (
	"Rabbit-OJ-Backend/services/channel"
	"context"
	"fmt"
	"sync"
)

var (
	MachineWaitGroup sync.WaitGroup
)

func StartMachine(ctx context.Context, index uint, queueChan chan []byte) {
	fmt.Printf("[Machine] Concurrent #%d started \n", index)
	MachineWaitGroup.Add(1)
	defer MachineWaitGroup.Done()

	for {
		select {
		case delivery := <-queueChan:
			fmt.Printf("[Machine] #%d machine START \n", index)
			okChan := make(chan bool)
			data := &channel.JudgeRequestBridgeMessage{
				Data:        delivery,
				SuccessChan: okChan,
			}
			channel.JudgeRequestBridgeChan <- data
			<-okChan
			close(okChan)
			fmt.Printf("[Machine] #%d machine FINISH \n", index)
		case <-ctx.Done():
			fmt.Printf("[Machine] #%d machine Exited \n", index)
			return
		}
	}
}