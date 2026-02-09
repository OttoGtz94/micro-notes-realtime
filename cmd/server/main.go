package main

import (
	"log"
	"net/http"
	// "realtime-service/handlers"
	"realtime-service/internal/ws"
)

func main() {
	// Registrar ruta
	// http.HandleFunc("/health", handlers.HealthHandler)
	// http.HandleFunc("/echo", handlers.EchoHandler)
	// http.HandleFunc("/events", handlers.EventsHandler)
	// http.HandleFunc("/ws", handlers.WSHandler)
	
	// V2
	// http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 	w.WriteHeader(http.StatusOK)
	// 	w.Write([]byte("OK"))
	// })

	// http.HandleFunc("/ws", ws.HandleWebSocket)

	// log.Println("Realtime service running on :8080")

	// err := http.ListenAndServe(":8080", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// V3
	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/ws", ws.HandleWebSocket(hub))

	log.Println("🚀 Go Realtime Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}