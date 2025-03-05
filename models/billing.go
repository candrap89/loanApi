package models

import "time"

type Billing struct {
	ID              int       `json:"id"`
	IDUser          int       `json:"id_user"`
	BillAmount      float64   `json:"bill_amount"`
	PaidStatus      bool      `json:"paid_status"`
	LastUpdatedAt   time.Time `json:"last_updated_at"`
	LoanOutstanding float64   `json:"loan_outstanding"`
	Week            int       `json:"week"`
}
