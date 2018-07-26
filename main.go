package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/robin1996/go-macro/mouse"
)

func main() {
	fmt.Println("start capturing mouse input")

	var isInterrupted bool
	var wg sync.WaitGroup

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	mouseChan := make(chan mouse.MouseMessage, 1)

	go func() {
		wg.Add(1)
		mouse.Notify(ctx, mouseChan)
		wg.Done()
	}()
	for {
		if isInterrupted {
			cancel()
			break
		}
		select {
		case <-signalChan:
			isInterrupted = true
		case k := <-mouseChan:
			fmt.Println(k.Button)
		}
	}
	wg.Wait()
	fmt.Println("done")
}
