package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/candrap89/loanApi/models"
	"github.com/candrap89/loanApi/queries"
	"github.com/segmentio/kafka-go"
)

// define globa variables for this file
var (
	ackReceived = make(map[string]bool)
	ackMutex    = &sync.Mutex{}
)

// class properties
type consumerHandler struct {
	UserLoanQuery queries.UserLoanQueryInterface // Use the interface
}

// NewConsumerHandler function to creates a new consumer handler with the given UserLoanQueryInterface
func NewConsumerHandler(userLoanQuery queries.UserLoanQueryInterface) *consumerHandler {
	return &consumerHandler{UserLoanQuery: userLoanQuery}
}

func (ch *consumerHandler) StartNewProductConsumer() {
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
		userLoan := models.UserLoan{
			UserCIF:         userCif,
			Loan:            5000000,
			Status:          true,
			LoanOutstanding: 5000000,
			IsDelinquent:    true,
		}

		err = ch.UserLoanQuery.CreateUserLoan(userLoan)
		if err != nil {
			fmt.Printf("Failed to insert User loan record: %v\n", err)
			return
		}
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
