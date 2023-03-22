package main

import (
	"fmt"
	"os"
)

func instrFinder(activeChan chan *input) {
	instrChanMap := make(map[string]chan *input)
	orderInstrMap := make(map[uint32]string)
	//reads from activeChan
	for {
		select {
		case inputPtr := <-activeChan:
			channel, prs := instrChanMap[inputPtr.instrument]
			//If not cancel orderCheck if instrument exists in map
			fmt.Fprintf(os.Stderr, "Instrument: %s\n", inputPtr.instrument)
			fmt.Fprintf(os.Stderr, "instrChanMap Length: %v\n", len(instrChanMap))
			if inputPtr.orderType == inputCancel {
				//guaranteed to have corresponding order to cancel somewhere due to discussion post
				//send cancel order correct chan
				instr := orderInstrMap[inputPtr.orderId]
				channel = instrChanMap[instr]
				channel <- inputPtr

			} else {
				if prs { //Yes: queue up active order to its handler
					orderInstrMap[inputPtr.orderId] = inputPtr.instrument //Todo: should this be before or after?
					channel <- inputPtr
				} else {
					fmt.Fprintf(os.Stderr, "Creating gmm\n")
					//Issue: buy order then cancel order
					//before gmm made for buyorder, cancel order try to find gmm for the same instr
					//dont have, so a second gmm for the same instr is made
					//
					//update instrChanMap to instantiate new handler
					newInstrChan := make(chan *input)
					instrChanMap[inputPtr.instrument] = newInstrChan
					go genericMatchmaker(newInstrChan)
					//queue up active order to its handler
					channel, prs := instrChanMap[inputPtr.instrument]
					orderInstrMap[inputPtr.orderId] = inputPtr.instrument //Todo: should this be before or after?
					channel <- inputPtr
					fmt.Fprintf(os.Stderr, "else prs: %v\n", prs)
				}
			}
		}
	}
}
