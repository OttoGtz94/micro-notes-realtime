package handlers

import (
	"io"
	"log"
	"net/http"
)

func EventsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// 👀 Aquí está la magia por ahora
	log.Println("📩 Event received:", string(body))

	w.WriteHeader(http.StatusAccepted)
}
