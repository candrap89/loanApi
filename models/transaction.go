package models

import "time"

// TransactionHistory represents a transaction record in the database
type TransactionHistory struct {
	ID        int       `json:"id"`         // Unique identifier for the transaction
	BillID    int       `json:"bill_id"`    // Foreign key referencing the billing table
	Amount    float64   `json:"amount"`     // Transaction amount
	CreatedAt time.Time `json:"created_at"` // Timestamp of the transaction
}
