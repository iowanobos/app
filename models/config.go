package models

type Config struct {
	RequestTopic  string `json:"requestTopic"`
	ResponseTopic string `json:"response_topic"`
	KafkaAddress  string `json:"kafkaAddress"`
}
