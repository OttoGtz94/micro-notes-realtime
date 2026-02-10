package main

import (
	"log"
	"net/http"
	httpHandlers "realtime-service/internal/http"
	"realtime-service/internal/ws"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()
	go ws.RedisListener(hub)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	
	http.HandleFunc("/ws", ws.HandleWebSocket(hub))

	http.HandleFunc("/events", httpHandlers.EventsHandler(hub))

	log.Println("🚀 Go Realtime Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}