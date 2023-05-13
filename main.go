package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	runServer(os.Args[1:])
}

func runServer(args []string) {
	log.SetFlags(0)

	var listenAddress = flag.String("listen", "127.0.0.1:1134", "Listen address.")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("HTTP %s %s%s\n", r.Method, r.Host, r.URL)

		if r.URL.Path != "/" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(time.Now().Local().Format(time.RFC1123Z)))
	})

	log.Printf("Listening at http://%s", *listenAddress)

	httpServer := http.Server{
		Addr: *listenAddress,
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-sigint
		fmt.Println("Shutting down...")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(idleConnectionsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	}

	<-idleConnectionsClosed

	fmt.Println("Bye bye")
}
