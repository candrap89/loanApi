package queries

import (
	"database/sql"

	"github.com/candrap89/loanApi/models"
)

type UserLoanQuery struct {
	DB *sql.DB
}

func NewUserLoanQuery(db *sql.DB) *UserLoanQuery {
	return &UserLoanQuery{DB: db}
}

func (q *UserLoanQuery) GetUserLoanByCIF(cif string) ([]models.UserLoan, error) {
	query := `
		SELECT id, user_cif, loan, status, last_updated_at, loan_outstanding, interest
		FROM user_loan
		WHERE user_cif = ?
	`

	rows, err := q.DB.Query(query, cif)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userLoans []models.UserLoan
	for rows.Next() {

		var userLoan models.UserLoan

		err := rows.Scan(
			&userLoan.ID,
			&userLoan.UserCIF,
			&userLoan.Loan,
			&userLoan.Status,
			&userLoan.LastUpdatedAt,
			&userLoan.LoanOutstanding,
			&userLoan.Interest,
		)
		if err != nil {
			return nil, err
		}

		userLoans = append(userLoans, userLoan)
	}

	return userLoans, nil
}

func (q *UserLoanQuery) GetAllUsers() ([]models.UserLoan, error) {
	query := `
		SELECT id, user_cif, loan, status, last_updated_at, loan_outstanding, interest
		FROM user_loan
	`

	rows, err := q.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userLoans []models.UserLoan
	for rows.Next() {
		var userLoan models.UserLoan
		err := rows.Scan(
			&userLoan.ID,
			&userLoan.UserCIF,
			&userLoan.Loan,
			&userLoan.Status,
			&userLoan.LastUpdatedAt,
			&userLoan.LoanOutstanding,
			&userLoan.Interest,
		)
		if err != nil {
			return nil, err
		}
		userLoans = append(userLoans, userLoan)
	}

	return userLoans, nil
}

// UpdateLoanOutstanding updates the loan_outstanding for a user
func (q *UserLoanQuery) UpdateUserLoanOutstanding(userID int, newOutstanding float64) error {
	query := `
		UPDATE user_loan
		SET loan_outstanding = ?
		WHERE id = ?
	`

	_, err := q.DB.Exec(query, newOutstanding, userID)
	return err
}

// Update user loan To deliquent updates the loan_outstanding for a user
func (q *UserLoanQuery) UpdateUserTodeliquent(IsDelinquent bool, user_id int) error {
	query := `
		UPDATE user_loan
		SET IsDelinquent = ?
		WHERE id = ?
	`

	_, err := q.DB.Exec(query, IsDelinquent, user_id)
	return err
}

func (q *UserLoanQuery) GetDelinquentUsers() ([]models.UserLoan, error) {
	query := `
		SELECT id, user_cif, loan, status, last_updated_at, loan_outstanding, interest, isDelinquent
		FROM user_loan where IsDelinquent = true
	`

	rows, err := q.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userLoans []models.UserLoan
	for rows.Next() {
		var userLoan models.UserLoan
		err := rows.Scan(
			&userLoan.ID,
			&userLoan.UserCIF,
			&userLoan.Loan,
			&userLoan.Status,
			&userLoan.LastUpdatedAt,
			&userLoan.LoanOutstanding,
			&userLoan.Interest,
			&userLoan.IsDelinquent,
		)
		if err != nil {
			return nil, err
		}
		userLoans = append(userLoans, userLoan)
	}

	return userLoans, nil
}
