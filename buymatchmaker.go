package main

import "fmt"

func buyMatchmaker(activeChan <-chan *order, matchedChan chan<- *order, unmatchedChan chan<- *order) {
	//takes in
	//reading end of activeChan
	//writing end of matchedChan
	//writing end of unmatchedChan
	for {
		select {
		case buyOrder := <-activeChan:
			match(buyOrder)
			matchedChan <- buyOrder   //Will block if unbuffered
			unmatchedChan <- buyOrder //Will block if unbuffered
		}
	}
}

func match(buyOrder *order) {
	fmt.Println("Matching...")
	printOrder(buyOrder)
}
