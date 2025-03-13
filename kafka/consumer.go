package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	ackReceived = make(map[string]bool)
	ackMutex    = &sync.Mutex{}
)

func StartNewProductConsumer() {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "new-products-ack",
		GroupID:  "inventory-service",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}

		var ack map[string]string
		json.Unmarshal(m.Value, &ack)

		userCif := ack["cifId"]
		ackMutex.Lock()
		ackReceived[userCif] = true
		ackMutex.Unlock()

		fmt.Printf("Received ACK for product: %s\n", userCif)
	}

	r.Close()
}

func WaitForAck(UserCif string) bool {
	for {
		ackMutex.Lock()
		if ackReceived[UserCif] {
			delete(ackReceived, UserCif)
			ackMutex.Unlock()
			return true
		}
		ackMutex.Unlock()
	}
}
