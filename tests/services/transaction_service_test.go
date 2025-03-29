package services_test

import (
	"errors"
	"go-json/internal/dtos/request"
	"go-json/internal/models"
	"go-json/internal/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) CreateTransaction(transaction models.Transaction) (*models.Transaction, error) {
	args := m.Called(transaction)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindAllTransaction() ([]models.Transaction, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Transaction), args.Error(1)
}

type TransactionServiceTestSuite struct {
	suite.Suite
	userRepo        *MockUserRepository
	transactionRepo *MockTransactionRepository
	roleRepo        *MockRoleRepository
	transactionSvc  services.TransactionService
	testUser        models.User
	testUserRoles   []models.UserRole
	testTransaction models.Transaction
}

func (suite *TransactionServiceTestSuite) SetupTest() {
	suite.userRepo = new(MockUserRepository)
	suite.transactionRepo = new(MockTransactionRepository)
	suite.roleRepo = new(MockRoleRepository)
	suite.transactionSvc = services.NewTransactionService(suite.userRepo, suite.transactionRepo, suite.roleRepo)

	suite.testUser = models.User{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Balance:  1000.0,
		IsActive: true,
	}

	suite.testUserRoles = []models.UserRole{
		{
			ID:     "1",
			UserID: "1",
			RoleID: "2",
		},
	}

	suite.testTransaction = models.Transaction{
		ID:           "1",
		CustomerID:   "1",
		MerchantID:   "2",
		ActivityType: models.PaymentActivity,
		Timestamp:    time.Now(),
		Details:      "Test transaction",
		Amount:       100.0,
	}
}

func (suite *TransactionServiceTestSuite) TestProcessPaymentSuccess() {
	paymentReq := request.PaymentRequest{
		CustomerID: "1",
		MerchantID: "2",
		Amount:     500.0,
	}

	updatedUser := suite.testUser
	updatedUser.Balance -= paymentReq.Amount
	suite.userRepo.On("FindByID", "1").Return(&suite.testUser, nil)
	suite.userRepo.On("UpdateUser", mock.AnythingOfType("models.User")).Return(nil)

	suite.roleRepo.On("FindRoleByUserID", "1").Return(&suite.testUserRoles, nil)

	expectedTransaction := models.Transaction{
		CustomerID:   "1",
		MerchantID:   "2",
		Amount:       500.0,
		ActivityType: models.PaymentActivity,
		Details:      "Payment processed successfully",
	}

	suite.transactionRepo.On("CreateTransaction", mock.MatchedBy(func(t models.Transaction) bool {
		return t.CustomerID == expectedTransaction.CustomerID &&
			t.MerchantID == expectedTransaction.MerchantID &&
			t.Amount == expectedTransaction.Amount &&
			t.ActivityType == expectedTransaction.ActivityType
	})).Return(&models.Transaction{
		ID:           "1",
		CustomerID:   "1",
		MerchantID:   "2",
		Amount:       500.0,
		ActivityType: models.PaymentActivity,
		Details:      "Payment processed successfully",
		Timestamp:    time.Now(),
	}, nil)

	response, err := suite.transactionSvc.ProcessPayment(paymentReq)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "1", response.ID)
	assert.Equal(suite.T(), "1", response.CustomerID)
	assert.Equal(suite.T(), "2", response.MerchantID)
	assert.Equal(suite.T(), 500.0, response.Amount)

	suite.userRepo.AssertExpectations(suite.T())
	suite.roleRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionServiceTestSuite) TestProcessPaymentInsufficientBalance() {
	paymentReq := request.PaymentRequest{
		CustomerID: "1",
		MerchantID: "2",
		Amount:     2000.0, // More than user's balance
	}

	// Setup user repo mock
	suite.userRepo.On("FindByID", "1").Return(&suite.testUser, nil)

	// Setup transaction repo mock - should record a failed payment
	suite.transactionRepo.On("CreateTransaction", mock.MatchedBy(func(t models.Transaction) bool {
		return t.CustomerID == "1" &&
			t.MerchantID == "2" &&
			t.Amount == 2000.0 &&
			t.ActivityType == models.FailedPayment &&
			t.Details == "Insufficient balance"
	})).Return(&models.Transaction{
		ID:           "1",
		CustomerID:   "1",
		MerchantID:   "2",
		Amount:       2000.0,
		ActivityType: models.FailedPayment,
		Details:      "Insufficient balance",
		Timestamp:    time.Now(),
	}, nil)

	// Test
	response, err := suite.transactionSvc.ProcessPayment(paymentReq)

	// Verify
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), "insufficient balance", err.Error())

	suite.userRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionServiceTestSuite) TestProcessPaymentInvalidCustomer() {
	paymentReq := request.PaymentRequest{
		CustomerID: "999", // Non-existent customer
		MerchantID: "2",
		Amount:     500.0,
	}

	// Setup user repo mock
	suite.userRepo.On("FindByID", "999").Return(nil, errors.New("user not found"))

	// Setup transaction repo mock - should record a failed payment
	suite.transactionRepo.On("CreateTransaction", mock.MatchedBy(func(t models.Transaction) bool {
		return t.CustomerID == "999" &&
			t.MerchantID == "2" &&
			t.Amount == 500.0 &&
			t.ActivityType == models.FailedPayment &&
			t.Details == "Invalid customer ID"
	})).Return(&models.Transaction{
		ID:           "1",
		CustomerID:   "999",
		MerchantID:   "2",
		Amount:       500.0,
		ActivityType: models.FailedPayment,
		Details:      "Invalid customer ID",
		Timestamp:    time.Now(),
	}, nil)

	// Test
	response, err := suite.transactionSvc.ProcessPayment(paymentReq)

	// Verify
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), "invalid customer ID", err.Error())

	suite.userRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionServiceTestSuite) TestProcessPaymentInactiveCustomer() {
	paymentReq := request.PaymentRequest{
		CustomerID: "1",
		MerchantID: "2",
		Amount:     500.0,
	}

	// Create inactive user
	inactiveUser := suite.testUser
	inactiveUser.IsActive = false

	// Setup user repo mock
	suite.userRepo.On("FindByID", "1").Return(&inactiveUser, nil)

	// Setup transaction repo mock - should record a failed payment
	suite.transactionRepo.On("CreateTransaction", mock.MatchedBy(func(t models.Transaction) bool {
		return t.CustomerID == "1" &&
			t.MerchantID == "2" &&
			t.Amount == 500.0 &&
			t.ActivityType == models.FailedPayment &&
			t.Details == "Customer is not active"
	})).Return(&models.Transaction{
		ID:           "1",
		CustomerID:   "1",
		MerchantID:   "2",
		Amount:       500.0,
		ActivityType: models.FailedPayment,
		Details:      "Customer is not active",
		Timestamp:    time.Now(),
	}, nil)

	// Test
	response, err := suite.transactionSvc.ProcessPayment(paymentReq)

	// Verify
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), "customer is not active", err.Error())

	suite.userRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionServiceTestSuite) TestProcessPaymentMerchantRole() {
	paymentReq := request.PaymentRequest{
		CustomerID: "1",
		MerchantID: "2",
		Amount:     500.0,
	}

	// Setup user with merchant role
	merchantRoles := []models.UserRole{
		{
			ID:     "1",
			UserID: "1",
			RoleID: "1", // merchant role
		},
	}

	// Setup user repo mock
	suite.userRepo.On("FindByID", "1").Return(&suite.testUser, nil)
	suite.userRepo.On("UpdateUser", mock.AnythingOfType("models.User")).Return(nil)

	// Setup role repo mock
	suite.roleRepo.On("FindRoleByUserID", "1").Return(&merchantRoles, nil)

	// Setup transaction repo mock
	suite.transactionRepo.On("CreateTransaction", mock.MatchedBy(func(t models.Transaction) bool {
		return t.CustomerID == "1" &&
			t.MerchantID == "2" &&
			t.Amount == 500.0 &&
			t.ActivityType == models.PaymentActivity
	})).Return(&models.Transaction{
		ID:           "1",
		CustomerID:   "1",
		MerchantID:   "2",
		Amount:       500.0,
		ActivityType: models.PaymentActivity,
		Details:      "Payment processed successfully",
		Timestamp:    time.Now(),
	}, nil)

	// Test
	response, err := suite.transactionSvc.ProcessPayment(paymentReq)

	// Verify
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	// Merchant should see both IDs
	assert.Equal(suite.T(), "2", response.MerchantID)

	suite.userRepo.AssertExpectations(suite.T())
	suite.roleRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionServiceTestSuite) TestTransactionHistoryByUserID() {
	transactions := []models.Transaction{
		suite.testTransaction,
		{
			ID:           "2",
			CustomerID:   "1",
			MerchantID:   "3",
			ActivityType: models.PaymentActivity,
			Timestamp:    time.Now(),
			Details:      "Another transaction",
			Amount:       200.0,
		},
		{
			ID:           "3",
			CustomerID:   "2", // Different customer
			MerchantID:   "3",
			ActivityType: models.PaymentActivity,
			Timestamp:    time.Now(),
			Details:      "Different customer",
			Amount:       300.0,
		},
	}

	// Setup mocks
	suite.userRepo.On("FindByID", "1").Return(&suite.testUser, nil)
	suite.transactionRepo.On("FindAllTransaction").Return(transactions, nil)

	// Test
	response, err := suite.transactionSvc.TransactionHistoryByUserID("1")

	// Verify
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Len(suite.T(), response, 1) // One user's history
	assert.Equal(suite.T(), "testuser", response[0].User.Username)
	assert.Len(suite.T(), response[0].Transactions, 2) // Two transactions for this user
	assert.Equal(suite.T(), 2, response[0].TotalCount)

	suite.userRepo.AssertExpectations(suite.T())
	suite.transactionRepo.AssertExpectations(suite.T())
}

func (suite *TransactionServiceTestSuite) TestTransactionHistoryUserNotFound() {
	// Setup mocks
	suite.userRepo.On("FindByID", "999").Return(nil, errors.New("user not found"))

	// Test
	response, err := suite.transactionSvc.TransactionHistoryByUserID("999")

	// Verify
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), "user not found", err.Error())

	suite.userRepo.AssertExpectations(suite.T())
}

func TestTransactionServiceSuite(t *testing.T) {
	suite.Run(t, new(TransactionServiceTestSuite))
}
