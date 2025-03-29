package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-json/internal/controllers"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) ProcessPayment(payment request.PaymentRequest) (*response.PaymentResponse, error) {
	args := m.Called(payment)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.PaymentResponse), args.Error(1)
}

func (m *MockTransactionService) TransactionHistoryByUserID(userID string) ([]response.UserTransactionHistoryResponse, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]response.UserTransactionHistoryResponse), args.Error(1)
}

type TransactionControllerTestSuite struct {
	suite.Suite
	transactionService *MockTransactionService
	controller         controllers.TransactionController
}

func (suite *TransactionControllerTestSuite) SetupTest() {
	suite.transactionService = new(MockTransactionService)
	suite.controller = controllers.NewTransactionController(suite.transactionService)
}

func (suite *TransactionControllerTestSuite) TestPayment() {
	paymentReq := request.PaymentRequest{
		CustomerID: "1",
		MerchantID: "2",
		Amount:     500.0,
	}

	paymentResp := &response.PaymentResponse{
		ID:         "1",
		CustomerID: "1",
		MerchantID: "2",
		Amount:     500.0,
		Timestamp:  time.Now(),
	}

	suite.transactionService.On("ProcessPayment", mock.MatchedBy(func(req request.PaymentRequest) bool {
		return req.CustomerID == paymentReq.CustomerID &&
			req.MerchantID == paymentReq.MerchantID &&
			req.Amount == paymentReq.Amount
	})).Return(paymentResp, nil)

	reqBody, _ := json.Marshal(paymentReq)
	req, _ := http.NewRequest("POST", "/payment", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Payment)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var apiResp response.ApiResponse
	err := json.Unmarshal(rr.Body.Bytes(), &apiResp)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, apiResp.Status)
	assert.Equal(suite.T(), "Payment successful", apiResp.Message)

	suite.transactionService.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestPaymentError() {
	paymentReq := request.PaymentRequest{
		CustomerID: "1",
		MerchantID: "2",
		Amount:     500.0,
	}

	suite.transactionService.On("ProcessPayment", mock.MatchedBy(func(req request.PaymentRequest) bool {
		return req.CustomerID == paymentReq.CustomerID &&
			req.MerchantID == paymentReq.MerchantID &&
			req.Amount == paymentReq.Amount
	})).Return(nil, errors.New("payment processing error"))

	reqBody, _ := json.Marshal(paymentReq)
	req, _ := http.NewRequest("POST", "/payment", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Payment)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)

	suite.transactionService.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestPaymentInvalidJSON() {
	reqBody := []byte(`{"bad json`)
	req, _ := http.NewRequest("POST", "/payment", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Payment)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func (suite *TransactionControllerTestSuite) TestTransactionHistory() {
	userID := "1"
	user := models.User{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
		Balance:  1000.0,
	}

	transactions := []*models.Transaction{
		{
			ID:           "1",
			CustomerID:   "1",
			MerchantID:   "2",
			ActivityType: models.PaymentActivity,
			Timestamp:    time.Now(),
			Details:      "Transaction 1",
			Amount:       100.0,
		},
		{
			ID:           "2",
			CustomerID:   "1",
			MerchantID:   "3",
			ActivityType: models.PaymentActivity,
			Timestamp:    time.Now(),
			Details:      "Transaction 2",
			Amount:       200.0,
		},
	}

	historyResponse := []response.UserTransactionHistoryResponse{
		{
			User:         user,
			Transactions: transactions,
			TotalCount:   len(transactions),
		},
	}

	suite.transactionService.On("TransactionHistoryByUserID", userID).Return(historyResponse, nil)

	// Create a new request
	req, _ := http.NewRequest("GET", "/transactions/1", nil)
	// Create a router to use the route variables
	router := mux.NewRouter()
	router.HandleFunc("/transactions/{id}", suite.controller.TransactionHistory)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	// Parse the response
	var apiResp response.ApiResponse
	err := json.Unmarshal(rr.Body.Bytes(), &apiResp)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, apiResp.Status)
	assert.Equal(suite.T(), "Transaction history retrieved", apiResp.Message)

	// Cast the data to the expected type and check values
	historyData, ok := apiResp.Data.([]interface{})
	assert.True(suite.T(), ok)
	assert.Len(suite.T(), historyData, 1)

	suite.transactionService.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestTransactionHistoryError() {
	userID := "999" // Non-existent user

	suite.transactionService.On("TransactionHistoryByUserID", userID).Return(nil, errors.New("user not found"))

	// Create a new request
	req, _ := http.NewRequest("GET", "/transactions/999", nil)
	// Create a router to use the route variables
	router := mux.NewRouter()
	router.HandleFunc("/transactions/{id}", suite.controller.TransactionHistory)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Serve the request
	router.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)

	suite.transactionService.AssertExpectations(suite.T())
}

func (suite *TransactionControllerTestSuite) TestTransactionHistoryNoID() {
	// Create a new request with no ID
	req, _ := http.NewRequest("GET", "/transactions/", nil)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Use the controller directly with an empty mux.Vars map
	req = mux.SetURLVars(req, map[string]string{})
	handler := http.HandlerFunc(suite.controller.TransactionHistory)
	handler.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(suite.T(), http.StatusBadRequest, rr.Code)
}

func TestTransactionControllerSuite(t *testing.T) {
	suite.Run(t, new(TransactionControllerTestSuite))
}
