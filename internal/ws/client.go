package ws

import (
	"log"
	"github.com/gorilla/websocket"
	"encoding/json"
	"time"
	"realtime-service/internal/events"
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
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