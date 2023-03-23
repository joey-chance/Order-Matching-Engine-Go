package main

import (
	"fmt"
	"os"
)

func instrFinder(activeChan chan input) {
	instrChanMap := make(map[string]chan order)
	orderInstrMap := make(map[uint32]string)
	//reads from activeChan
	var currTime int64 = 0
	for {
		select {
		case inputItem := <-activeChan:
			channel, prs := instrChanMap[inputItem.instrument]

			// Get timestamp here
			//var timestamp int64 = GetCurrentTimestamp()

			orderItem := order{&inputItem, currTime, 1, 0}
			currTime++

			//If not cancel orderCheck if instrument exists in map
			fmt.Fprintf(os.Stderr, "Instrument: %s\n", inputItem.instrument)
			fmt.Fprintf(os.Stderr, "instrChanMap Length: %v\n", len(instrChanMap))
			if inputItem.orderType == inputCancel {
				//guaranteed to have corresponding order to cancel somewhere due to discussion post
				//send cancel order correct chan
				instr := orderInstrMap[inputItem.orderId]
				channel = instrChanMap[instr]
				channel <- orderItem

			} else {
				if prs { //Yes: queue up active order to its handler
					orderInstrMap[inputItem.orderId] = inputItem.instrument //Todo: should this be before or after?
					channel <- orderItem
				} else {
					fmt.Fprintf(os.Stderr, "Creating gmm\n")
					//Issue: buy order then cancel order
					//before gmm made for buyorder, cancel order try to find gmm for the same instr
					//dont have, so a second gmm for the same instr is made
					//
					//update instrChanMap to instantiate new handler
					newInstrChan := make(chan order, 100)
					instrChanMap[inputItem.instrument] = newInstrChan
					go genericMatchmaker(newInstrChan)
					//queue up active order to its handler
					channel, prs := instrChanMap[inputItem.instrument]
					orderInstrMap[inputItem.orderId] = inputItem.instrument //Todo: should this be before or after?
					channel <- orderItem
					fmt.Fprintf(os.Stderr, "else prs: %v\n", prs)
				}
			}
		}
	}
}
