package services

import (
	"awesomeProject1/diploma/models"
	"awesomeProject1/diploma/models/request"
	"awesomeProject1/diploma/utils"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/google/uuid"
	"log"
	"time"
)

type Sender struct {
	Name         string
	sendTopic    string
	receiveTopic string
	producer     sarama.SyncProducer
	consumer     sarama.ConsumerGroup
}

func (s Sender) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s Sender) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (s Sender) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Println("Got a response message: " + string(msg.Value))
	}
	return nil
}

func NewSender(c *models.Config) *Sender {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{c.KafkaAddress}, config)
	if err != nil {
		log.Fatal("Create producer failed. Error: " + err.Error())
	}
	consumer, err := sarama.NewConsumerGroup([]string{c.KafkaAddress}, "sender-group", config)
	if err != nil {
		log.Fatal("Create consumer failed. Error: " + err.Error())
	}

	return &Sender{
		Name:         "sender",
		sendTopic:    c.RequestTopic,
		receiveTopic: c.ResponseTopic,
		producer:     producer,
		consumer:     consumer,
	}
}

func (s *Sender) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second).C
	go s.consumer.Consume(ctx, []string{s.receiveTopic}, s)
	for {
		select {
		case <-ticker:
			message := s.generateRequestMessage()
			rawMessage, err := json.Marshal(message)
			if err != nil {
				log.Println("Serialize the request message failed. Error: " + err.Error())
				continue
			}

			partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
				Topic: s.sendTopic,
				Value: sarama.ByteEncoder(rawMessage),
			})
			if err != nil {
				log.Println("Send the request message failed. Error: " + err.Error())
				break
			}
			log.Println("Send the message succeeded. Value: "+string(rawMessage)+". Partition: ", partition, ". Offset: ", offset)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Sender) generateRequestMessage() *request.MessageModel {
	return &request.MessageModel{
		Header: request.MessageHeaderModel{
			Title:         utils.GetRandText() + "-request title",
			ID:            uuid.New().String(),
			Timestamp:     time.Now(),
			Service:       s.Name,
			CorrelationID: "null",
		},
		Body: request.MessageBodyModel{
			Method: utils.GetRandText() + "-method",
			Params: nil,
		},
	}
}
