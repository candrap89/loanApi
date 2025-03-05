package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/candrap89/loanApi/queries"
)

type UserLoanHandler struct {
	UserLoanQuery *queries.UserLoanQuery
}

func NewUserLoanHandler(userLoanQuery *queries.UserLoanQuery) *UserLoanHandler {
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
