package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		sign := <-stop
		fmt.Printf("awaiting signal: %v", sign)
		wg.Done()
	}()
	fmt.Println("awaiting signal")
	wg.Wait()
}
