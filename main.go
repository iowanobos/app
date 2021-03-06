package main

import (
	"github.com/iowanobos/app/clients"
	"github.com/iowanobos/app/models"
	"github.com/iowanobos/app/services"

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

	consul := clients.NewConsul()

	receiver := services.NewReceiver(config)
	go receiver.Start(ctx)
	consul.Register(receiver.Name)

	sender := services.NewSender(config)
	go sender.Start(ctx)
	consul.Register(sender.Name)

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
		consul.Unregister(receiver.Name)
		consul.Unregister(sender.Name)
		<-ctx.Done()
	}
	log.Println("Application shut downing...")
}
