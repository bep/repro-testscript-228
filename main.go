package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	runServer(os.Args[1:])
}

func runServer(args []string) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sigint
	fmt.Println("Shutting down...")

}
