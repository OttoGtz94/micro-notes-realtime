package handlers

import (
	"io"
	"net/http"
	"fmt"
)

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	// Leer todo el body de la request
	body, err := io.ReadAll(r.Body)
	fmt.Println("Body:", string(body))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Devolver exactamente lo mismo
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}