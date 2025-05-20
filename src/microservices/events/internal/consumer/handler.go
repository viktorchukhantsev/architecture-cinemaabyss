package consumer

import (
	"encoding/json"
	"log"

	"github.com/Shopify/sarama"

	"events/internal/consts"
	"events/internal/dto"
)

func (s *SimpleConsumer) RunHandleMessages() {
	go s.handleMessages(consts.MoviesEventsTopis)
	go s.handleMessages(consts.UsersEventsTopis)
	go s.handleMessages(consts.PaymentsEventsTopis)
}

func (s *SimpleConsumer) handleMessages(topic string) {
	partitionConsumer, err := s.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Printf("Failed to create partition consumer for topic %s: %v", topic, err)
		return
	}
	defer partitionConsumer.Close()

	log.Printf("Started consuming messages from topic: %s", topic)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Received message from topic %s: %s", topic, string(msg.Value))
			s.processMessage(topic, msg.Value)
		case err := <-partitionConsumer.Errors():
			log.Printf("Error consuming from topic %s: %v", topic, err)
		}
	}
}

func (s *SimpleConsumer) processMessage(topic string, message []byte) {
	switch topic {
	case consts.MoviesEventsTopis:
		var event dto.Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Error unmarshaling movie event: %v", err)
			return
		}
		log.Printf("Processing movie event: %+v", event)
	case consts.UsersEventsTopis:
		var event dto.Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Error unmarshaling user event: %v", err)
			return
		}
		log.Printf("Processing user event: %+v", event)
	case consts.PaymentsEventsTopis:
		var event dto.Event
		if err := json.Unmarshal(message, &event); err != nil {
			log.Printf("Error unmarshaling payment event: %v", err)
			return
		}
		log.Printf("Processing payment event: %+v", event)
	}
}
