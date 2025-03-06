package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/candrap89/loanApi/mocks"
	"github.com/candrap89/loanApi/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPaymentHandler_MakePayment_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockBillingQuery := mocks.NewMockBillingQuery(ctrl)
	mockUserLoanQuery := mocks.NewMockUserLoanQuery(ctrl)
	mockTransactionQuery := mocks.NewMockTransactionQuery(ctrl)

	// Create the handler with mocks
	handler := NewPaymentHandler(mockBillingQuery, mockUserLoanQuery, mockTransactionQuery)

	// Define test data
	userId := 1
	amount := 2000.00
	billingRecords := []models.Billing{
		{
			ID:              1,
			IDUser:          userId,
			BillAmount:      1000.00,
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: 5500000.00,
			Week:            1,
		},
		{
			ID:              2,
			IDUser:          userId,
			BillAmount:      1000.00,
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: 5500000.00,
			Week:            2,
		},
	}

	// Set up mock expectations
	mockBillingQuery.EXPECT().GetBillingByID(userId).Return(billingRecords, nil)
	mockUserLoanQuery.EXPECT().UpdateUserLoanOutstanding(userId, billingRecords[1].LoanOutstanding-amount).Return(nil)
	mockBillingQuery.EXPECT().MarkBillAsPaidandUpdateOutstanding(billingRecords[1].LoanOutstanding-amount, billingRecords[1].ID).Return(nil)
	mockBillingQuery.EXPECT().MarkBillAsPaid(billingRecords[0].ID).Return(nil)
	mockBillingQuery.EXPECT().MarkBillAsPaid(billingRecords[1].ID).Return(nil)
	mockUserLoanQuery.EXPECT().UpdateUserTodeliquent(false, userId).Return(nil)
	mockTransactionQuery.EXPECT().InsertTransaction(gomock.Any()).Return(nil)

	// Create a request
	requestBody, _ := json.Marshal(PaymentRequest{UserId: userId, Amount: amount})
	req, _ := http.NewRequest(http.MethodPost, "/payment", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.MakePayment(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rr.Code)

	var response PaymentResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "Payment successful", response.Message)
}

func TestPaymentHandler_MakePayment_InvalidRequestBody(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockBillingQuery := mocks.NewMockBillingQuery(ctrl)
	mockUserLoanQuery := mocks.NewMockUserLoanQuery(ctrl)
	mockTransactionQuery := mocks.NewMockTransactionQuery(ctrl)

	// Create the handler with mocks
	handler := NewPaymentHandler(mockBillingQuery, mockUserLoanQuery, mockTransactionQuery)

	// Create a request with invalid body
	req, _ := http.NewRequest(http.MethodPost, "/payment", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.MakePayment(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "Invalid request body")
}

func TestPaymentHandler_MakePayment_FailedToFetchBilling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockBillingQuery := mocks.NewMockBillingQuery(ctrl)
	mockUserLoanQuery := mocks.NewMockUserLoanQuery(ctrl)
	mockTransactionQuery := mocks.NewMockTransactionQuery(ctrl)

	// Create the handler with mocks
	handler := NewPaymentHandler(mockBillingQuery, mockUserLoanQuery, mockTransactionQuery)

	// Define test data
	userId := 1
	amount := 2000.00

	// Set up mock expectations
	mockBillingQuery.EXPECT().GetBillingByID(userId).Return(nil, errors.New("database error"))

	// Create a request
	requestBody, _ := json.Marshal(PaymentRequest{UserId: userId, Amount: amount})
	req, _ := http.NewRequest(http.MethodPost, "/payment", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.MakePayment(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to fetch billing record")
}

func TestPaymentHandler_MakePayment_PaymentAmountMismatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockBillingQuery := mocks.NewMockBillingQuery(ctrl)
	mockUserLoanQuery := mocks.NewMockUserLoanQuery(ctrl)
	mockTransactionQuery := mocks.NewMockTransactionQuery(ctrl)

	// Create the handler with mocks
	handler := NewPaymentHandler(mockBillingQuery, mockUserLoanQuery, mockTransactionQuery)

	// Define test data
	userId := 1
	amount := 1000.00 // Incorrect amount
	billingRecords := []models.Billing{
		{
			ID:              1,
			IDUser:          userId,
			BillAmount:      1000.00,
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: 5500000.00,
			Week:            1,
		},
		{
			ID:              2,
			IDUser:          userId,
			BillAmount:      1000.00,
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: 5500000.00,
			Week:            2,
		},
	}

	// Set up mock expectations
	mockBillingQuery.EXPECT().GetBillingByID(userId).Return(billingRecords, nil)

	// Create a request
	requestBody, _ := json.Marshal(PaymentRequest{UserId: userId, Amount: amount})
	req, _ := http.NewRequest(http.MethodPost, "/payment", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.MakePayment(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response PaymentResponse
	json.NewDecoder(rr.Body).Decode(&response)
	assert.Equal(t, http.StatusBadRequest, response.Code)
	assert.Equal(t, "Payment amount does not match the bill amount", response.Message)
}

func TestPaymentHandler_MakePayment_FailedToUpdateLoanOutstanding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mocks
	mockBillingQuery := mocks.NewMockBillingQuery(ctrl)
	mockUserLoanQuery := mocks.NewMockUserLoanQuery(ctrl)
	mockTransactionQuery := mocks.NewMockTransactionQuery(ctrl)

	// Create the handler with mocks
	handler := NewPaymentHandler(mockBillingQuery, mockUserLoanQuery, mockTransactionQuery)

	// Define test data
	userId := 1
	amount := 2000.00
	billingRecords := []models.Billing{
		{
			ID:              1,
			IDUser:          userId,
			BillAmount:      1000.00,
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: 5500000.00,
			Week:            1,
		},
		{
			ID:              2,
			IDUser:          userId,
			BillAmount:      1000.00,
			PaidStatus:      false,
			LastUpdatedAt:   time.Now(),
			LoanOutstanding: 5500000.00,
			Week:            2,
		},
	}

	// Set up mock expectations
	mockBillingQuery.EXPECT().GetBillingByID(userId).Return(billingRecords, nil)
	mockUserLoanQuery.EXPECT().UpdateUserLoanOutstanding(userId, billingRecords[1].LoanOutstanding-amount).Return(errors.New("database error"))

	// Create a request
	requestBody, _ := json.Marshal(PaymentRequest{UserId: userId, Amount: amount})
	req, _ := http.NewRequest(http.MethodPost, "/payment", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the handler
	handler.MakePayment(rr, req)

	// Assert the response
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "Failed to update loan outstanding")
}
