package models

import "time"

type UserLoan struct {
	ID              int       `json:"id"`
	UserCIF         string    `json:"user_cif"`
	Loan            float64   `json:"loan"`
	Status          bool      `json:"status"`
	LastUpdatedAt   time.Time `json:"last_updated_at"`
	LoanOutstanding float64   `json:"loan_outstanding"`
	Interest        float64   `json:"interest"`
	IsDelinquent    bool      `json:"is_delinquent"`
}
