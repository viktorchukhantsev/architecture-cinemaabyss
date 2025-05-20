package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"events/internal/consumer"
	"events/internal/dto"
	"events/internal/producer"
)

var (
	simpleProducer *producer.SimpleProducer
	simpleConsumer *consumer.SimpleConsumer
)

func main() {
	kafkaBroker := os.Getenv("KAFKA_BROKERS")

	simpleProducer = producer.MustNewSimpleSyncProducer(kafkaBroker)
	defer simpleProducer.Close()

	simpleConsumer = consumer.MustNewSimpleConsumer(kafkaBroker)
	defer simpleConsumer.Close()
	simpleConsumer.RunHandleMessages()

	router := mux.NewRouter()
	router.HandleFunc("/api/events/health", handleHealth).Methods("GET")
	router.HandleFunc("/api/events/movie", handleMovieEvent).Methods("POST")
	router.HandleFunc("/api/events/user", handleUserEvent).Methods("POST")
	router.HandleFunc("/api/events/payment", handlePaymentEvent).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	log.Printf("Starting events service on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"status": true})
}

func handleMovieEvent(w http.ResponseWriter, r *http.Request) {
	var movieEvent dto.MovieEvent
	if err := json.NewDecoder(r.Body).Decode(&movieEvent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	partition, offset, event, err := simpleProducer.ProduceMovieEvent(movieEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"partition": partition,
		"offset":    offset,
		"event":     event,
	})
}

func handleUserEvent(w http.ResponseWriter, r *http.Request) {
	var userEvent dto.UserEvent
	if err := json.NewDecoder(r.Body).Decode(&userEvent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	partition, offset, event, err := simpleProducer.ProduceUserEvent(userEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"partition": partition,
		"offset":    offset,
		"event":     event,
	})
}

func handlePaymentEvent(w http.ResponseWriter, r *http.Request) {
	var paymentEvent dto.PaymentEvent
	if err := json.NewDecoder(r.Body).Decode(&paymentEvent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	partition, offset, event, err := simpleProducer.ProducePaymentEvent(paymentEvent)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "success",
		"partition": partition,
		"offset":    offset,
		"event":     event,
	})
}
