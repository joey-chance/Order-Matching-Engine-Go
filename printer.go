package main

import (
	"fmt"
	"os"
)

type printObj struct {
	outputType   string //Exec, Delete, Add
	restingOrder order
	activeOrder  order
	execCount    uint32
	execPrice    uint32
	timestamp    int64
	canCancel    bool
}

func printer(printChan <-chan printObj) {
	for printJob := range printChan {
		// fmt.Fprintf(os.Stderr, "printer: %v\n", printJob)
		// fmt.Fprintf(os.Stderr, "activeorder: %v|%v\n", printJob.activeOrder.inp.orderId, printJob.outputType)
		if printJob.outputType == "E" {
			outputOrderExecuted(printJob.restingOrder.inp.orderId,
				printJob.activeOrder.inp.orderId,
				printJob.restingOrder.executionId,
				printJob.execPrice,
				printJob.execCount,
				printJob.timestamp)
		} else if printJob.outputType == "A" {
			outputOrderAdded(*printJob.activeOrder.inp, printJob.timestamp)
		} else if printJob.outputType == "D" {
			outputOrderDeleted(*printJob.activeOrder.inp, printJob.canCancel, printJob.timestamp)
		} else {
			fmt.Fprintf(os.Stderr, "dropped order: %v|%v\n", printJob.activeOrder.inp.orderId, printJob.outputType)
		}
	}
}
