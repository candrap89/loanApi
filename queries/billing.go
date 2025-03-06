package queries

import (
	"database/sql"

	"github.com/candrap89/loanApi/models"
)

type BillingQuery struct {
	DB *sql.DB
}

func NewBillingQuery(db *sql.DB) *BillingQuery {
	return &BillingQuery{DB: db}
}

// InsertBilling inserts a new billing record into the database
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

// GetLatestWeek
func (q *BillingQuery) GetLatestWeek(userID int) (int, error) {
	query := `
		SELECT week
		FROM billing
		WHERE id_user = ?
		ORDER BY week DESC
		LIMIT 1
	`

	var week int
	err := q.DB.QueryRow(query, userID).Scan(&week)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no records exist for the user, start from week 1
			return 1, nil
		}
		return 0, err
	}

	return week + 1, nil // Increment the latest week by 1
}

// getDeliquentUsers fetches all users who have outstanding loan
func (q *BillingQuery) GetDeliquentUsers() ([]int, error) {
	query := `
		SELECT id_user
		FROM billing
		WHERE paid_status = false
		GROUP BY id_user
		HAVING COUNT(id_user) > 1
	`

	rows, err := q.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var delinquentUsers []int
	for rows.Next() {
		var userID int
		err := rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		delinquentUsers = append(delinquentUsers, userID)
	}

	return delinquentUsers, nil
}

// GetBillingByID fetches a billing record by ID
func (q *BillingQuery) GetBillingByID(userID int) ([]models.Billing, error) {
	query := `
		SELECT id, id_user, bill_amount, paid_status, last_updated_at, loan_outstanding, week
		FROM billing
		WHERE id_user = ? AND paid_status = false
	`

	rows, err := q.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var billings []models.Billing
	for rows.Next() {
		var billing models.Billing
		err := rows.Scan(
			&billing.ID,
			&billing.IDUser,
			&billing.BillAmount,
			&billing.PaidStatus,
			&billing.LastUpdatedAt,
			&billing.LoanOutstanding,
			&billing.Week,
		)
		if err != nil {
			return nil, err
		}
		billings = append(billings, billing)
	}

	return billings, nil
}

// MarkBillAsPaidandUpdateOutstanding updates the paid_status of a billing record to true and updates the loan_outstanding
func (q *BillingQuery) MarkBillAsPaidandUpdateOutstanding(outstanding float64, billID int) error {
	query := `
		UPDATE billing
		SET paid_status = true,
		loan_outstanding = ?
		WHERE id = ?
	`
	_, err := q.DB.Exec(query, outstanding, billID)
	return err
}

// MarkBillAsPaid updates the paid_status of a billing record to true
func (q *BillingQuery) MarkBillAsPaid(billID int) error {
	query := `
		UPDATE billing
		SET paid_status = true
		WHERE id = ?
	`
	_, err := q.DB.Exec(query, billID)
	return err
}
