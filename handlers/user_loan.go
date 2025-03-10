package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/candrap89/loanApi/queries"
)

type UserLoanHandler struct {
	UserLoanQuery queries.UserLoanQueryInterface // Use the interface
}

func NewUserLoanHandler(userLoanQuery queries.UserLoanQueryInterface) *UserLoanHandler {
	return &UserLoanHandler{UserLoanQuery: userLoanQuery}
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
