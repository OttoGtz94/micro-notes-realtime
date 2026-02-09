package http

import (
	"encoding/json"
	"log"
	nethttp "net/http"
	"time"

	"realtime-service/internal/events"
	"realtime-service/internal/ws"
)

func EventsHandler(hub *ws.Hub) nethttp.HandlerFunc {
	return func(w nethttp.ResponseWriter, r *nethttp.Request) {

		if r.Method != nethttp.MethodPost {
			w.WriteHeader(nethttp.StatusMethodNotAllowed)
			return
		}

		var incoming struct {
			Type    events.EventType `json:"type"`
			Payload any              `json:"payload"`
		}

		err := json.NewDecoder(r.Body).Decode(&incoming)
		if err != nil {
			w.WriteHeader(nethttp.StatusBadRequest)
			return
		}

		event := events.Event{
			Type:      incoming.Type,
			Payload:   incoming.Payload,
			Timestamp: time.Now().UTC(),
		}

		bytes, err := json.Marshal(event)
		if err != nil {
			w.WriteHeader(nethttp.StatusInternalServerError)
			return
		}

		log.Printf("📨 HTTP Event: %s", bytes)

		hub.Broadcast(bytes)

		w.WriteHeader(nethttp.StatusAccepted)
	}
}
