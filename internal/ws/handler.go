package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade error:", err)
			return
		}

		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte),
		}

		client.hub.register <- client

		go client.writePump()
		go client.readPump()
	}
}

// package ws

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gorilla/websocket"
// )

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // por ahora, luego lo aseguramos
// 	},
// }

// func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("WS upgrade error:", err)
// 		return
// 	}
// 	defer conn.Close()

// 	log.Println("🟢 Cliente conectado por WebSocket")

// 	for {
// 		messageType, message, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println("Cliente desconectado")
// 			break
// 		}

// 		log.Printf("📩 Mensaje recibido: %s\n", message)

// 		// Echo (responde lo mismo)
// 		err = conn.WriteMessage(messageType, message)
// 		if err != nil {
// 			log.Println("Error al escribir mensaje:", err)
// 			break
// 		}
// 	}
// }
