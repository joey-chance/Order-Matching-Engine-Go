package main

import "fmt"

func sellMatchmaker(activeChan <-chan *order, matchedChan chan<- *order, unmatchedChan chan<- *order) {
	//takes in
	//reading end of activeChan
	//writing end of matchedChan
	//writing end of unmatchedChan
	for {
		select {
		case sellOrder := <-activeChan:
			match(sellOrder)
			matchedChan <- sellOrder   //Will block if unbuffered
			unmatchedChan <- sellOrder //Will block if unbuffered
		}
	}
}

func match(sellOrder *order) {
	fmt.Println("Matching...")
	printOrder(sellOrder)
}
