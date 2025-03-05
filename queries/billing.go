package queries

import (
	"database/sql"

	"github.com/loanApi/models"
)

type BillingQuery struct {
	DB *sql.DB
}

func NewBillingQuery(db *sql.DB) *BillingQuery {
	return &BillingQuery{DB: db}
}

func (q *BillingQuery) InsertBilling(billing models.Billing) error {
	query := `
		INSERT INTO billing (id_user, bill_amount, paid_status, last_updated_at, loan_outstanding, week)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := q.DB.Exec(
		query,
		billing.IDUser,
		billing.BillAmount,
		billing.PaidStatus,
		billing.LastUpdatedAt,
		billing.LoanOutstanding,
		billing.Week,
	)

	return err
}
