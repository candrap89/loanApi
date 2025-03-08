package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/candrap89/loanApi/mocks"
	"github.com/candrap89/loanApi/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserLoanByCIF(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserLoanQuery := mocks.NewMockUserLoanQuery(ctrl)
	handler := NewUserLoanHandler(mockUserLoanQuery)

	// Define test cases
	tests := []struct {
		name           string
		cif            string
		mockReturn     []models.UserLoan
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			cif:  "12345",
			mockReturn: []models.UserLoan{
				{
					ID:              1,
					UserCIF:         "12345",
					Loan:            1000.0,
					Status:          true,
					LastUpdatedAt:   time.Now(),
					LoanOutstanding: 500.0,
					Interest:        10.0,
					IsDelinquent:    false,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_cif":"12345","loan":1000,"status":true,"last_updated_at":"`,
		},
		{
			name:           "CIF Parameter Missing",
			cif:            "",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `CIF parameter is required`,
		},
		{
			name:           "User Loan Not Found",
			cif:            "12345",
			mockReturn:     []models.UserLoan{},
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `User loan data not found`,
		},
		{
			name:           "Database Error",
			cif:            "12345",
			mockReturn:     nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cif != "" {
				mockUserLoanQuery.EXPECT().
					GetUserLoanByCIF(tt.cif).
					Return(tt.mockReturn, tt.mockError).
					Times(1)
			}

			req, err := http.NewRequest("GET", "/user-loan?cif="+tt.cif, nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.GetUserLoanByCIF(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestGetDelinquentUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserLoanQuery := mocks.NewMockUserLoanQuery(ctrl)
	handler := NewUserLoanHandler(mockUserLoanQuery)

	// Define test cases
	tests := []struct {
		name           string
		mockReturn     []models.UserLoan
		mockError      error
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			mockReturn: []models.UserLoan{
				{
					ID:              1,
					UserCIF:         "12345",
					Loan:            1000.0,
					Status:          true,
					LastUpdatedAt:   time.Now(),
					LoanOutstanding: 500.0,
					Interest:        10.0,
					IsDelinquent:    true,
				},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"user_cif":"12345","loan":1000,"status":true,"last_updated_at":"`,
		},
		{
			name:           "No Delinquent Users",
			mockReturn:     []models.UserLoan{},
			mockError:      nil,
			expectedStatus: http.StatusNotFound,
			expectedBody:   `Deliquent User loan data not found`,
		},
		{
			name:           "Database Error",
			mockReturn:     nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserLoanQuery.EXPECT().
				GetDelinquentUsers().
				Return(tt.mockReturn, tt.mockError).
				Times(1)

			req, err := http.NewRequest("GET", "/delinquent-users", nil)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.GetDelinquentUsers(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}
