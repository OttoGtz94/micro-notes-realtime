package ws

import (
	"encoding/json"
	"log"
	"realtime-service/internal/events"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			var event events.Event

			if err := json.Unmarshal(message, &event); err != nil {
				log.Println("Error unmarshalling event:", err)
				continue
			}

			for client := range h.clients {
				if event.SenderID != "" && client.id == event.SenderID {
					log.Println("🚫 Filtrando eco para:", client.id)
					continue
				}
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}

func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}
