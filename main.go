package main

import (
	"awesomeProject1/diploma/models"
	"awesomeProject1/diploma/services"
	"context"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	config := &models.Config{
		RequestTopic:  "request-topic",
		ResponseTopic: "response-topic",
		KafkaAddress:  "localhost:9092",
	}

	receiver := services.NewReceiver(config)
	go receiver.Start(ctx)
	sender := services.NewSender(config)
	go sender.Start(ctx)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	select {
	case <-ctx.Done():
	case <-sigc:
		cancel()
		<-ctx.Done()
	}
	log.Println("Application shut downing...")
}
