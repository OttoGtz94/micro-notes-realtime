package events

import "time"

type Event struct {
	Module    string    `json:"module"`
	Type      EventType `json:"type"`
	Payload   any       `json:"payload"`
	SenderID  string    `json:"senderId"`
	Timestamp time.Time `json:"timestamp"`
}
