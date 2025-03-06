package queries

import (
	"database/sql"

	"github.com/candrap89/loanApi/models"
)

type TransactionQuery struct {
	DB *sql.DB
}

func NewTransactionQuery(db *sql.DB) *TransactionQuery {
	return &TransactionQuery{DB: db}
}

// InsertTransaction inserts a new transaction record into the database
func (q *TransactionQuery) InsertTransaction(transaction models.TransactionHistory) error {
	query := `
		INSERT INTO transaction_history (bill_id, amount, created_at)
		VALUES (?, ?, ?)
	`

	_, err := q.DB.Exec(
		query,
		transaction.BillID,
		transaction.Amount,
		transaction.CreatedAt,
	)

	return err
}
