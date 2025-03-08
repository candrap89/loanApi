package queries

import "github.com/candrap89/loanApi/models"

// BillingQueryInterface defines the methods for BillingQuery
type BillingQueryInterface interface {
	GetBillingByID(userID int) ([]models.Billing, error)
	MarkBillAsPaidandUpdateOutstanding(outstanding float64, billID int) error
	MarkBillAsPaid(billID int) error
}

// UserLoanQueryInterface defines the methods for UserLoanQuery
type UserLoanQueryInterface interface {
	UpdateUserLoanOutstanding(userID int, newOutstanding float64) error
	UpdateUserTodeliquent(status bool, userID int) error
	GetDelinquentUsers() ([]models.UserLoan, error)
	GetUserLoanByCIF(cif string) ([]models.UserLoan, error)
}

// TransactionQueryInterface defines the methods for TransactionQuery
type TransactionQueryInterface interface {
	InsertTransaction(transaction models.TransactionHistory) error
}
