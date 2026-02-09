package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Esto "convierte" HTTP en WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // solo para desarrollo
	},
}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WS upgrade error:", err)
		return
	}
	// Cerrar la conexión cuando la función termine
	defer conn.Close()

	log.Println("🔌 WebSocket client connected")

	// Mensaje inicial
	conn.WriteMessage(
		websocket.TextMessage,
		[]byte("welcome from go 👋"),
	)

	// 🔁 Loop infinito: mantiene viva la conexión
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("🔴 WebSocket client disconnected")
			break
		}

		log.Printf("📨 Message received: %s\n", message)

		// (opcional) eco al cliente
		conn.WriteMessage(messageType, message)
	}
}
