package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/candrap89/loanApi/models"
	"github.com/candrap89/loanApi/queries"
	"github.com/segmentio/kafka-go"
)

type produceMessage struct {
	UserLoanQuery queries.UserLoanQueryInterface
}

func NewProductMessage(userLoanQuery queries.UserLoanQueryInterface) *produceMessage {
	return &produceMessage{UserLoanQuery: userLoanQuery}
}

func (pm *produceMessage) SendNewProductMessage(barcode string) error {
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "new-products", 0)
	if err != nil {
		fmt.Printf("Failed to connect to Kafka leader: %v\n", err)
		return err
	}
	defer conn.Close()

	message := map[string]string{
		"barcode": barcode,
	}

	// create userloan simul
	userLoan := models.UserLoan{
		UserCIF:         "12345TESTkafka",
		Loan:            5000000,
		Status:          true,
		LoanOutstanding: 5000000,
		IsDelinquent:    true,
	}

	// Assuming ch is an instance of a struct with UserLoanQuery as a field

	err = pm.UserLoanQuery.CreateUserLoan(userLoan)
	if err != nil {
		fmt.Printf("Failed to insert User loan record: %v\n", err)
		return err
	}

	msg, _ := json.Marshal(message)
	conn.WriteMessages(
		kafka.Message{Value: msg},
	)
	return err
}
