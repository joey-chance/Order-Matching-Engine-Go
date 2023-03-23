package main

//Generic matchmaker for both buys & sells (instr level conc)

import (
	"container/heap"
	"fmt"
	"os"
)

func genericMatchmaker(activeInstrChan <-chan order) {
	bpq := make(BuyPriorityQueue, 0)
	spq := make(SellPriorityQueue, 0)
	heap.Init(&bpq)
	heap.Init(&spq)
	iter := 0
	for {
		select {
		case orderItem := <-activeInstrChan:
			fmt.Fprintf(os.Stderr, "iter: %v\n", iter)
			iter++

			orderptr := &orderItem

			switch orderptr.inp.orderType {
			case inputCancel:
				{
					fmt.Fprintf(os.Stderr, "CANCEL: \n")

					// if instr != orderptr.inp.instrument { //dro
					// 	fmt.Fprintf(os.Stderr, "My Instr: %s\n", instr)
					// 	fmt.Fprintf(os.Stderr, "Order Instr: %s\n", orderptr.inp.instrument)
					// 	continue
					// }
					orderFound := false
					// check if in any pq
					//fmt.Fprintf(os.Stderr, "bpq new len: %v\n", len(bpq))
					removedOrderPtr := (&bpq).RemoveOrderId(orderptr.inp.orderId)

					if removedOrderPtr != nil {
						fmt.Fprintf(os.Stderr, "bpq removedOrderPtr: %v\n", removedOrderPtr)
						go outputOrderDeleted(*orderptr.inp, true, orderptr.timestamp)
						continue
					}

					removedOrderPtr = (&spq).RemoveOrderId(orderptr.inp.orderId)
					if removedOrderPtr != nil {
						fmt.Fprintf(os.Stderr, "spq removedOrderPtr: %v\n", removedOrderPtr)
						go outputOrderDeleted(*orderptr.inp, true, orderptr.timestamp)
						continue
					}
					// no output false
					if !orderFound {
						fmt.Fprintf(os.Stderr, "rej removedOrderPtr: %v\n", removedOrderPtr)
						go outputOrderDeleted(*orderptr.inp, false, orderptr.timestamp)
					}
				}
			case inputBuy:
				{
					//Try to match as many as possible in opposing pq
					buyOrderPtr := orderptr
					//while spq has items, activeOrder count>0,
					for spq.Len() > 0 && buyOrderPtr.inp.count > 0 {
						//peek top

						if spq.Peek().(*order).inp.price <= buyOrderPtr.inp.price {
							sellOrder := heap.Pop(&spq).(*order)
							execPrice := sellOrder.inp.price //Always execute at lower price, which is sell price

							if buyOrderPtr.inp.count < sellOrder.inp.count { //if sell order partial match, edit sellOrder,increment executionId, decrement activeOrder count, outputOrderExecuted
								execCount := buyOrderPtr.inp.count
								go outputOrderExecuted(sellOrder.inp.orderId, buyOrderPtr.inp.orderId, sellOrder.executionId, execPrice, execCount, buyOrderPtr.timestamp)
								//modify execution id of sell order
								sellOrder.inp.count += 1
								//modify count of sell order
								sellOrder.inp.count -= execCount
								//put sellorder back into spq
								heap.Push(&spq, sellOrder)
								//set buyorder to 0
								buyOrderPtr.inp.count = 0
							} else if buyOrderPtr.inp.count == sellOrder.inp.count { //if sell order full match, pop, outputOrderExecuted, decrement activeOrder count and try again if still have remaining count
								execCount := buyOrderPtr.inp.count
								go outputOrderExecuted(sellOrder.inp.orderId, buyOrderPtr.inp.orderId, sellOrder.executionId, execPrice, execCount, buyOrderPtr.timestamp)
								//set buyorder to 0
								buyOrderPtr.inp.count = 0
							} else if buyOrderPtr.inp.count > sellOrder.inp.count {
								execCount := sellOrder.inp.count
								go outputOrderExecuted(sellOrder.inp.orderId, buyOrderPtr.inp.orderId, sellOrder.executionId, execPrice, execCount, buyOrderPtr.timestamp)
								//set buyorder to 0
								buyOrderPtr.inp.count -= execCount
							}
						} else {
							break
						}
					}
					fmt.Fprintf(os.Stderr, "spq new len: %v\n", len(spq))

					//Finally Check if activeOrder has remaining qty
					//Yes push activeOrder to pq, outputOrderAdded
					//no do nothing
					if buyOrderPtr.inp.count > 0 {
						heap.Push(&bpq, buyOrderPtr)
						fmt.Fprintf(os.Stderr, "bpq new len: %v\n", len(bpq))
						go outputOrderAdded(*buyOrderPtr.inp, buyOrderPtr.timestamp)
					}
				}
			case inputSell:
				{
					//Try to match as many as possible in opposing pq
					sellOrderPtr := orderptr
					//while bpq has items, activeOrder count>0,
					for bpq.Len() > 0 && sellOrderPtr.inp.count > 0 {
						//peek top
						fmt.Fprintf(os.Stderr, "bpq beset Price: %v\n", bpq.Peek().(*order).inp.price)
						fmt.Fprintf(os.Stderr, "bpq Length: %v\n", bpq.Len())
						for idx, item := range bpq {
							fmt.Fprintf(os.Stderr, "bpq item: %v|%v\n", item.inp.price, idx)
						}
						if bpq.Peek().(*order).inp.price >= sellOrderPtr.inp.price {
							buyOrder := heap.Pop(&bpq).(*order)
							execPrice := sellOrderPtr.inp.price //Always execute at lower price, which is sell price

							if sellOrderPtr.inp.count < buyOrder.inp.count { //if sell order partial match, edit buyOrder,increment executionId, decrement activeOrder count, outputOrderExecuted
								execCount := sellOrderPtr.inp.count
								go outputOrderExecuted(buyOrder.inp.orderId, sellOrderPtr.inp.orderId, buyOrder.executionId, execPrice, execCount, sellOrderPtr.timestamp)
								//modify execution id of sell order
								buyOrder.inp.count += 1
								//modify count of sell order
								buyOrder.inp.count -= execCount
								//put buyorder back into bpq
								heap.Push(&bpq, buyOrder)
								//set buyorder to 0
								sellOrderPtr.inp.count = 0
							} else if sellOrderPtr.inp.count == buyOrder.inp.count { //if sell order full match, pop, outputOrderExecuted, decrement activeOrder count and try again if still have remaining count
								execCount := sellOrderPtr.inp.count
								go outputOrderExecuted(buyOrder.inp.orderId, sellOrderPtr.inp.orderId, buyOrder.executionId, execPrice, execCount, sellOrderPtr.timestamp)
								//set buyorder to 0
								sellOrderPtr.inp.count = 0
							} else if sellOrderPtr.inp.count > buyOrder.inp.count {
								execCount := buyOrder.inp.count
								go outputOrderExecuted(buyOrder.inp.orderId, sellOrderPtr.inp.orderId, buyOrder.executionId, execPrice, execCount, sellOrderPtr.timestamp)
								//set buyorder to 0
								sellOrderPtr.inp.count -= execCount
							}
						} else {
							break
						}
					}

					//Finally Check if activeOrder has remaining qty
					//Yes push activeOrder to pq, outputOrderAdded
					//no do nothing
					if sellOrderPtr.inp.count > 0 {
						heap.Push(&spq, sellOrderPtr)
						go outputOrderAdded(*sellOrderPtr.inp, sellOrderPtr.timestamp)
					}
				}
			}
		}
	}
}
