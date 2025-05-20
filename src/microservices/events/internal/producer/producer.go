package producer

import (
	"log"

	"github.com/Shopify/sarama"
)

type SimpleProducer struct {
	config   *sarama.Config
	producer sarama.SyncProducer
}

func MustNewSimpleSyncProducer(broker string) *SimpleProducer {
	var res SimpleProducer

	var newBroker string

	if broker == "" {
		newBroker = "kafka:9092"
	} else {
		newBroker = broker
	}

	brokers := []string{newBroker}

	res.config = sarama.NewConfig()
	res.config.Producer.RequiredAcks = sarama.WaitForAll
	res.config.Producer.Retry.Max = 5
	res.config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, res.config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	res.producer = producer

	log.Println("Kafka producer initialized successfully")

	return &res
}

func (c *SimpleProducer) Close() {
	_ = c.producer.Close()
}
