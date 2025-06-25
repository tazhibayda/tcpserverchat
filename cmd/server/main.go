package main

import (
	"context"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"tcpserverchat/internal/server"
)

func main() {
	srv := server.New(":8000")
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("failed to start: %v", err)
		}
	}()
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":9090", nil)
	}()
	go func() {
		log.Println("pprof on :6060")
		http.ListenAndServe(":6060", nil)
	}()
	<-ctx.Done()
	log.Println("Shutting down...")
	srv.Stop()
	log.Println("Stopped")
}
