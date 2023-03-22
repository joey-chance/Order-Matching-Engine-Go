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
	fmt.Println("Hello, world!")
	//Buyspq Testing
	order0 := order{&input{'B', 1, 5, 100, "AAPL"}, 0, 0}
	order1 := order{&input{'B', 1, 10, 100, "AAPL"}, 1, 0}
	order2 := order{&input{'B', 1, 3, 100, "AAPL"}, 2, 0}
	order3 := order{&input{'B', 1, 5, 100, "AAPL"}, 3, 0}

	//init empty pq
	pq := make(BuyPriorityQueue, 0)
	heap.Init(&pq)
	//push orders in 1 by 1
	heap.Push(&pq, &order1)
	heap.Push(&pq, &order2)
	heap.Push(&pq, &order3)
	heap.Push(&pq, &order0)
	//pop orders out, should be in order: order 1, 0, 3, 2
	for pq.Len() > 0 {
		poppedorder := heap.Pop(&pq).(*order) //type assertion to order type
		var poppedinput = poppedorder.inp
		fmt.Println(poppedorder)
		fmt.Println(poppedinput)
		fmt.Println("")
	}

	//Sellspq Testing
	sorder0 := order{&input{'S', 1, 5, 100, "AAPL"}, 0, 0}
	sorder1 := order{&input{'S', 1, 10, 100, "AAPL"}, 1, 0}
	sorder2 := order{&input{'S', 1, 3, 100, "AAPL"}, 2, 0}
	sorder3 := order{&input{'S', 1, 5, 100, "AAPL"}, 3, 0}

	//init empty pq
	spq := make(SellPriorityQueue, 0)
	heap.Init(&spq)
	//push orders in 1 by 1
	heap.Push(&spq, &sorder1)
	heap.Push(&spq, &sorder2)
	heap.Push(&spq, &sorder3)
	heap.Push(&spq, &sorder0)
	//pop orders out, should be in order: order 2, 0 ,3, 1
	for spq.Len() > 0 {
		poppedorder := heap.Pop(&spq).(*order) //type assertion to order type
		var poppedinput = poppedorder.inp
		fmt.Println(poppedorder)
		fmt.Println(poppedinput)
		fmt.Println("")
	}
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
