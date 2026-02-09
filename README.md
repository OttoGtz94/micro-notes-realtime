# Go Realtime Service — Fase Broadcast + Eventos JSON

> Documentación técnica y pedagógica del microservicio **Realtime en Go**

Este documento explica **qué se construyó**, **por qué**, **qué archivos siguen siendo útiles**, cuáles **ya no se usan**, y cómo funciona todo el sistema **paso a paso**, pensado para alguien que **está iniciando en Go**.

---

## 1. Contexto general

Este microservicio tiene un solo propósito:

> 📡 **Recibir eventos del backend (NestJS en el futuro) y transmitirlos en tiempo real al frontend usando WebSockets**

Principios clave:

- Go **solo maneja realtime**
- No hay lógica de negocio
- No hay base de datos
- Todo gira alrededor de **eventos**

---

## 2. Estado actual del árbol de directorios

```txt
.
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
├── handlers              ❌ LEGACY (ya no se usan)
│   ├── echo.go
│   ├── events.go
│   ├── health.go
│   └── ws.go
├── hub                   ❌ LEGACY (primer intento)
│   ├── client.go
│   └── hub.go
└── internal               ✅ CÓDIGO REAL
    ├── events
    │   ├── event.go
    │   └── types.go
    └── ws
        ├── client.go
        ├── handler.go
        └── hub.go
```

---

## 3. ¿Qué carpetas siguen sirviendo y cuáles no?

### ✅ Carpetas ACTIVAS

#### `cmd/server`

- Punto de entrada del programa
- Arranca el servidor HTTP
- Inicializa el Hub

#### `internal/ws`

Contiene **toda la infraestructura realtime**:

- conexiones WebSocket
- manejo de clientes
- broadcast

#### `internal/events`

Define el **contrato de eventos JSON**:

- tipos permitidos
- estructura estándar

---

### ❌ Carpetas LEGACY (pueden eliminarse)

#### `handlers/`

- Era la versión inicial
- HTTP handlers simples
- Ya **no se usan**

#### `hub/`

- Primer experimento del Hub
- Reemplazado por `internal/ws/hub.go`

👉 **Recomendación**: puedes borrarlas sin miedo

```bash
rm -rf handlers hub
```

---

## 4. Arquitectura final (conceptual)

```txt
              ┌─────────────┐
              │   NestJS     │
              │ (futuro)     │
              └──────┬──────┘
                     │ HTTP /events
                     ▼
              ┌─────────────┐
              │   Go API     │
              │  (this svc)  │
              └──────┬──────┘
                     │ broadcast
                     ▼
              ┌─────────────┐
              │     HUB     │
              └──────┬──────┘
                     │ WebSocket
        ┌────────────┼────────────┐
        ▼            ▼            ▼
   Frontend A   Frontend B   Frontend C
```

---

## 5. Flujo completo de un evento

1. Un cliente se conecta vía `/ws`
2. Se registra en el Hub
3. Un mensaje entra (temporalmente desde WS)
4. Se transforma en **Evento JSON**
5. El Hub lo envía a **todos los clientes conectados**

---

## 6. Código explicado (archivo por archivo)

---

### `cmd/server/main.go`

```go
package main

import (
	"log"
	"net/http"

	"realtime-service/internal/ws"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/ws", ws.HandleWebSocket(hub))

	log.Println("🚀 Go Realtime Service running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

🧠 Qué hace:

- Crea el Hub
- Lo ejecuta en una goroutine
- Expone `/health` y `/ws`

---

### `internal/ws/hub.go`

```go
package ws

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
			delete(h.clients, client)
			close(client.send)

		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}
```

🧠 Qué es el Hub:

- Controlador central
- No entiende eventos
- Solo distribuye bytes

---

### `internal/ws/client.go`

```go
package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
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
			Type:      events.NoteUpdated,
			Payload:   string(message),
			Timestamp: time.Now().UTC(),
		}

		bytes, _ := json.Marshal(event)
		log.Printf("📦 Event: %s", bytes)

		c.hub.broadcast <- bytes
	}
}

func (c *Client) writePump() {
	for msg := range c.send {
		c.conn.WriteMessage(1, msg)
	}
}
```

🧠 Responsabilidad:

- Leer mensajes
- Convertirlos en eventos
- Enviarlos al Hub

---

### `internal/ws/handler.go`

```go
package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebSocket(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}

		client := &Client{hub: hub, conn: conn, send: make(chan []byte)}
		hub.register <- client

		go client.readPump()
		go client.writePump()
	}
}
```

---

### `internal/events/types.go`

```go
package events

type EventType string

const (
	NoteUpdated    EventType = "NOTE_UPDATED"
	NoteRestored   EventType = "NOTE_RESTORED"
	SyncConflict   EventType = "SYNC_CONFLICT"
	LoginFailed   EventType = "LOGIN_FAILED"
	TokenRefreshed EventType = "TOKEN_REFRESHED"
)
```

---

### `internal/events/event.go`

```go
package events

import "time"

type Event struct {
	Type      EventType `json:"type"`
	Payload   any       `json:"payload"`
	Timestamp time.Time `json:"timestamp"`
}
```

---

## 7. Estado actual del proyecto

✅ Broadcast funcional
✅ Eventos JSON tipados
✅ Arquitectura limpia
✅ Listo para recibir eventos HTTP desde NestJS

---

## 8. Próximo paso recomendado

➡️ **Endpoint HTTP `/events`**

NestJS → Go → WebSocket

Este será el punto donde ambos mundos se conectan.

---
