package consumer

import (
	"log"

	"github.com/Shopify/sarama"
)

type SimpleConsumer struct {
	config   *sarama.Config
	consumer sarama.Consumer
}

func MustNewSimpleConsumer(broker string) *SimpleConsumer {
	var res SimpleConsumer

	var newBroker string

	if broker == "" {
		newBroker = "kafka:9092"
	} else {
		newBroker = broker
	}

	brokers := []string{newBroker}

	res.config = sarama.NewConfig()
	res.config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer(brokers, res.config)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	res.consumer = consumer

	log.Println("Kafka consumer initialized successfully")

	return &res
}

func (c *SimpleConsumer) Close() {
	_ = c.consumer.Close()
}
