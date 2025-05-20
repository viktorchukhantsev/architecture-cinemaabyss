package producer

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Shopify/sarama"

	"events/internal/consts"
	"events/internal/dto"
)

func (s *SimpleProducer) ProduceMovieEvent(movieEvent dto.MovieEvent) (int32, int64, dto.Event, error) {
	event := dto.Event{
		ID:        eventID(dto.MovieEventType, movieEvent.MovieID, movieEvent.Action),
		Type:      dto.MovieEventType,
		Timestamp: time.Now(),
		Payload:   movieEvent,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return 0, 0, event, err
	}

	msg := &sarama.ProducerMessage{
		Topic: consts.MoviesEventsTopis,
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, event, err
	}

	log.Printf("Movie event sent to partition %d at offset %d", partition, offset)

	return partition, offset, event, nil
}

func (s *SimpleProducer) ProduceUserEvent(userEvent dto.UserEvent) (int32, int64, dto.Event, error) {
	event := dto.Event{
		ID:        eventID(dto.UserEventType, userEvent.UserID, userEvent.Action),
		Type:      dto.UserEventType,
		Timestamp: time.Now(),
		Payload:   userEvent,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return 0, 0, event, err
	}

	msg := &sarama.ProducerMessage{
		Topic: consts.UsersEventsTopis,
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, event, err
	}

	log.Printf("User event sent to partition %d at offset %d", partition, offset)

	return partition, offset, event, nil
}

func (s *SimpleProducer) ProducePaymentEvent(paymentEvent dto.PaymentEvent) (int32, int64, dto.Event, error) {
	event := dto.Event{
		ID:        eventID(dto.PaymentEventType, paymentEvent.PaymentID, paymentEvent.Status),
		Type:      dto.PaymentEventType,
		Timestamp: time.Now(),
		Payload:   paymentEvent,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return 0, 0, event, err
	}

	// Send event to Kafka
	msg := &sarama.ProducerMessage{
		Topic: consts.PaymentsEventsTopis,
		Value: sarama.StringEncoder(eventJSON),
	}

	partition, offset, err := s.producer.SendMessage(msg)
	if err != nil {
		return 0, 0, event, err
	}

	log.Printf("Payment event sent to partition %d at offset %d", partition, offset)

	return partition, offset, event, nil
}

func eventID(eventType dto.EventType, id int, action string) string {
	return fmt.Sprintf("%s-%d-%s", eventType, id, action)
}
