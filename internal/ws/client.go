package ws

import (
	"encoding/json"
	"log"
	"realtime-service/internal/events"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	id   string
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		event := events.Event{
			Type:      events.NoteUpdated, // por ahora fijo
			Payload:   string(message),
			Timestamp: time.Now().UTC(),
		}

		bytes, err := json.Marshal(event)
		if err != nil {
			log.Println("JSON error:", err)
			continue
		}

		log.Printf("📦 Event broadcast: %s", bytes)
		c.hub.broadcast <- bytes
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for message := range c.send {
		err := c.conn.WriteMessage(1, message)
		if err != nil {
			break
		}
	}
}