package main

import "C"
import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type Engine struct {
}

var activeChan chan input

func init() {
	activeChan = make(chan input)
	go instrFinder(activeChan)
}

func (e *Engine) accept(ctx context.Context, conn net.Conn) {
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	go handleConn(conn, activeChan)
}

func handleConn(conn net.Conn, activeChan chan<- input) {
	defer conn.Close()
	for {
		in, err := readInput(conn)
		if err != nil {
			if err != io.EOF {
				_, _ = fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			}
			return
		}
		fmt.Fprintf(os.Stderr, "Reading input\n")
		activeChan <- in
		fmt.Fprintf(os.Stderr, "Finished reading input\n")
	}
}

func GetCurrentTimestamp() int64 {
	return time.Now().UnixNano()
}
