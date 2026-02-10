package ws

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func RedisListener(hub *Hub) {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     "192.168.28.174:6379",
		Password: "P4ssw0rd",
		DB:       0,
	})

	pubsub := rdb.Subscribe(ctx, "app_event_bus")
	defer pubsub.Close()

	log.Println("📻 Escuchando el bus de eventos global 'app_event_bus'...")

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("📥 Mensaje recibido de Redis: %s", msg.Payload)

		hub.broadcast <- []byte(msg.Payload)
	}
}
