package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/candrap89/loanApi/kafka"
	"github.com/candrap89/loanApi/models"
	"github.com/candrap89/loanApi/queries"
)

type UserLoanHandler struct {
	UserLoanQuery queries.UserLoanQueryInterface // Use the interface
}

func NewUserLoanHandler(userLoanQuery queries.UserLoanQueryInterface) *UserLoanHandler {
	return &UserLoanHandler{UserLoanQuery: userLoanQuery}
}

type CreateLoanResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *UserLoanHandler) CreateUserLoan(w http.ResponseWriter, r *http.Request) {
	var req models.UserLoan
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userLoan := models.UserLoan{
		UserCIF:         req.UserCIF,
		Loan:            req.Loan,
		Status:          req.Status,
		LoanOutstanding: req.LoanOutstanding,
		Interest:        req.Interest,
		IsDelinquent:    req.IsDelinquent,
	}

	err := h.UserLoanQuery.CreateUserLoan(userLoan)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert User loan record: %v", err), http.StatusInternalServerError)
		return
	}
	// Return success response
	response := PaymentResponse{
		Code:    http.StatusOK,
		Message: "User Loan created successfully",
	}

	// Send a Kafka message to notify other departments
	kafka.SendNewProductMessage(userLoan.UserCIF)

	// Wait for ACK from the consumer
	if !kafka.WaitForAck(userLoan.UserCIF) {
		http.Error(w, "Timeout waiting for ACK", http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func (h *UserLoanHandler) GetUserLoanByCIF(w http.ResponseWriter, r *http.Request) {
	cif := r.URL.Query().Get("cif")
	if cif == "" {
		http.Error(w, "CIF parameter is required", http.StatusBadRequest)
		return
	}

	userLoans, err := h.UserLoanQuery.GetUserLoanByCIF(cif)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(userLoans) == 0 {
		http.Error(w, "User loan data not found", http.StatusNotFound)
		return
	}
	fmt.Println("User loans:", userLoans)
	ctx := r.Context()
	fmt.Println("User loan data found in context:", ctx.Value("user_loan_data"))

	ctx = context.WithValue(ctx, "user_loan_data :", userLoans)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userLoans)

}

func (h *UserLoanHandler) GetDelinquentUsers(w http.ResponseWriter, r *http.Request) {
	userLoans, err := h.UserLoanQuery.GetDelinquentUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(userLoans) == 0 {
		http.Error(w, "Deliquent User loan data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userLoans)

}
