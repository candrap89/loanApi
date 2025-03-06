package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/candrap89/loanApi/models"
	"github.com/candrap89/loanApi/queries"
)

type PaymentHandler struct {
	BillingQuery     queries.BillingQueryInterface
	UserLoanQuery    queries.UserLoanQueryInterface
	TransactionQuery queries.TransactionQueryInterface
}

func NewPaymentHandler(billingQuery queries.BillingQueryInterface, userLoanQuery queries.UserLoanQueryInterface, transactionQuery queries.TransactionQueryInterface) *PaymentHandler {
	return &PaymentHandler{
		BillingQuery:     billingQuery,
		UserLoanQuery:    userLoanQuery,
		TransactionQuery: transactionQuery,
	}
}

// PaymentRequest represents the input for the payment API
type PaymentRequest struct {
	UserId int     `json:"user_id"`
	Amount float64 `json:"amount"`
}

// PaymentResponse represents the JSON response for the payment API
type PaymentResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// MakePayment handles the payment API
func (h *PaymentHandler) MakePayment(w http.ResponseWriter, r *http.Request) {
	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Fetch the billing record
	billing, err := h.BillingQuery.GetBillingByID(req.UserId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch billing record: %v", err), http.StatusInternalServerError)
		return
	}

	// get toiotal bill amount
	var billAmount float64
	var lastbill = models.Billing{}
	for _, bill := range billing {
		billAmount = bill.BillAmount + billAmount
		lastbill = bill
	}

	// Validate the payment amount
	if req.Amount != billAmount {

		response := PaymentResponse{
			Code:    http.StatusBadRequest,
			Message: "Payment amount does not match the bill amount",
		}

		// Return failed response
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Update the loan_outstanding in the user_loan table based on latest billing record
	err = h.UserLoanQuery.UpdateUserLoanOutstanding(lastbill.IDUser, lastbill.LoanOutstanding-req.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update loan outstanding: %v", err), http.StatusInternalServerError)
		return
	}

	// Update the paid_status in the billing table based on latest billing record
	err = h.BillingQuery.MarkBillAsPaidandUpdateOutstanding(lastbill.LoanOutstanding-req.Amount, lastbill.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update billing record: %v", err), http.StatusInternalServerError)
		return
	}

	// if evrything ok update billing record as paid
	for _, bill := range billing {
		err = h.BillingQuery.MarkBillAsPaid(bill.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update billing record: %v", err), http.StatusInternalServerError)
			return
		}

	}
	// updaye user to not deliquent
	err = h.UserLoanQuery.UpdateUserTodeliquent(false, lastbill.IDUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user record: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert a transaction record if evrytihng is successful
	transaction := models.TransactionHistory{
		BillID:    lastbill.ID,
		Amount:    req.Amount,
		CreatedAt: time.Now(),
	}
	err = h.TransactionQuery.InsertTransaction(transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert transaction record: %v", err), http.StatusInternalServerError)
		return
	}
	// Return success response
	response := PaymentResponse{
		Code:    http.StatusOK,
		Message: "Payment successful",
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
