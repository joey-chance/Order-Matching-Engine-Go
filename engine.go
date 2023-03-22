package main

import "C"
import (
	"container/heap"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func init() {
	//Engine channels
	activeChan := make(chan *order)
	unmatchedChan := make(chan *order)
	matchedChan := make(chan *order)

	go buyMatchmaker(activeChan, matchedChan, unmatchedChan)

	bpq := make(BuyPriorityQueue, 0)
	spq := make(SellPriorityQueue, 0)
	heap.Init(&bpq)
	heap.Init(&spq)

	activeChan <- &order{&input{'B', 1, 5, 100, "AAPL"}, 0, 0}

	//Receives from channels to unblock
	fmt.Println("Matched Channel received:")
	printOrder(<-matchedChan)
	fmt.Println("Unmatched Channel received:")
	printOrder(<-unmatchedChan)
}

type Engine struct{}

func (e *Engine) accept(ctx context.Context, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	go handleConn(conn)
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		in, err := readInput(conn)
		if err != nil {
			if err != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}
			return
		}
		switch in.orderType {
		case inputCancel:
			fmt.Fprintf(os.Stderr, "Got cancel ID: %v\n", in.orderId)
			outputOrderDeleted(in, true, GetCurrentTimestamp())
		default:
			fmt.Fprintf(os.Stderr, "Got order: %c %v x %v @ %v ID: %v\n",
				in.orderType, in.instrument, in.count, in.price, in.orderId)
			outputOrderAdded(in, GetCurrentTimestamp())
		}
		outputOrderExecuted(123, 124, 1, 2000, 10, GetCurrentTimestamp())
	}
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano()
}
