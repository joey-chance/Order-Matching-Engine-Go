package main

import "fmt"

func printOrder(order *order) {
	fmt.Println("|", "Instr:", order.inp.instrument, "|",
		"Count:", order.inp.count, "|",
		"Price:", order.inp.price, "|",
		"Time:", order.timestamp, "|",
		"Type:", string(order.inp.orderType), "|",
		"Id", order.inp.orderId, "|",
		"Index", order.index, "|")
}
