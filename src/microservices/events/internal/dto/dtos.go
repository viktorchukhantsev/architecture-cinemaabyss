package dto

import "time"

type EventType = string

const (
	PaymentEventType EventType = "payment"
	UserEventType    EventType = "user"
	MovieEventType   EventType = "movie"
)

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}

type MovieEvent struct {
	MovieID     int      `json:"movie_id"`
	Title       string   `json:"title"`
	Action      string   `json:"action"` // viewed, rated, added, etc.
	UserID      int      `json:"user_id,omitempty"`
	Rating      float64  `json:"rating,omitempty"`
	Genres      []string `json:"genres,omitempty"`
	Description string   `json:"description,omitempty"`
}

type UserEvent struct {
	UserID    int       `json:"user_id"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Action    string    `json:"action"` // registered, logged_in, updated_profile, etc.
	Timestamp time.Time `json:"timestamp"`
}

type PaymentEvent struct {
	PaymentID  int       `json:"payment_id"`
	UserID     int       `json:"user_id"`
	Amount     float64   `json:"amount"`
	Status     string    `json:"status"` // completed, failed, refunded, etc.
	Timestamp  time.Time `json:"timestamp"`
	MethodType string    `json:"method_type,omitempty"`
}
