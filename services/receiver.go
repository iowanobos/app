package services

import (
	"awesomeProject1/diploma/models"
	"awesomeProject1/diploma/models/request"
	"awesomeProject1/diploma/models/response"
	"awesomeProject1/diploma/utils"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"math/rand"
	"time"
)

type Receiver struct {
	Name         string
	sendTopic    string
	receiveTopic string
	producer     sarama.SyncProducer
	consumer     sarama.ConsumerGroup
}

func (r Receiver) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (r Receiver) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (r Receiver) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Println("Got a request message: " + string(msg.Value))
		var req request.MessageModel
		err := json.Unmarshal(msg.Value, &req)
		if err != nil {
			log.Println("Deserialize request message failed. Error: " + err.Error() + " Message: " + string(msg.Value))
			continue
		}

		message := r.generateResponseMessage(req.Header.ID)
		messageJSON, err := json.Marshal(message)
		if err != nil {
			log.Println("Serialize response message failed. Error: " + err.Error())
			continue
		}
		r.producer.SendMessage(&sarama.ProducerMessage{
			Topic: r.sendTopic,
			Value: sarama.ByteEncoder(messageJSON),
		})
	}
	return nil
}

func NewReceiver(c *models.Config) *Receiver {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{c.KafkaAddress}, config)
	if err != nil {
		log.Fatal("Create producer failed. Error: " + err.Error())
	}
	consumer, err := sarama.NewConsumerGroup([]string{c.KafkaAddress}, "receiver-group", config)
	if err != nil {
		log.Fatal("Create consumer failed. Error: " + err.Error())
	}

	return  &Receiver{
		Name:         "receiver",
		sendTopic:    c.ResponseTopic,
		receiveTopic: c.RequestTopic,
		producer:     producer,
		consumer:     consumer,
	}
}

func (r *Receiver) Start(ctx context.Context) {
	go r.consumer.Consume(ctx, []string{r.receiveTopic}, r)

}

func (r *Receiver) generateResponseMessage(id string) *response.MessageModel {
	var messageBody response.MessageBodyModel
	if rand.Intn(100) > 20 {
		result := utils.GetRandText() + "-result"
		messageBody = response.MessageBodyModel{
			Success: true,
			Result:  &result,
		}
	} else {
		err := response.ErrorModel{
			ErrorCode:    "500",
			ErrorMessage: "INTERNAL_SERVER_ERROR",
		}
		messageBody = response.MessageBodyModel{
			Success: false,
			Error:   &err,
		}
	}

	return &response.MessageModel{
		Header: response.MessageHeaderModel{
			Title:         utils.GetRandText() + "-response title",
			ID:            id,
			Timestamp:     time.Now(),
			Service:       r.Name,
			CorrelationID: "null",
			Status:        utils.GetRandText() + "-status",
		},
		Body: messageBody,
	}
}
