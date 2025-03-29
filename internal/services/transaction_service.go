package services

import (
	"errors"
	"go-json/internal/dtos/mapper"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/models"
	"go-json/internal/repositories"
	"time"

	"github.com/go-playground/validator/v10"
)

type TransactionService interface {
	ProcessPayment(payment request.PaymentRequest) (*response.PaymentResponse, error)
	TransactionHistoryByUserID(userID string) ([]response.UserTransactionHistoryResponse, error)
}

type transactionService struct {
	userRepo        repositories.UserRepository
	transactionRepo repositories.TransactionRepository
	roleRepo        repositories.RoleRepository
}

func NewTransactionService(userRepo repositories.UserRepository, transactionRepo repositories.TransactionRepository, roleRepo repositories.RoleRepository) TransactionService {
	return &transactionService{userRepo: userRepo, transactionRepo: transactionRepo, roleRepo: roleRepo}
}

func (p *transactionService) ProcessPayment(payment request.PaymentRequest) (*response.PaymentResponse, error) {
	validate := validator.New()
	err := validate.Struct(payment)
	if err != nil {
		return nil, err
	}
	var transaction models.Transaction
	transaction.CustomerID = payment.CustomerID
	transaction.MerchantID = payment.MerchantID
	transaction.Amount = payment.Amount
	transaction.Timestamp = time.Now()
	transaction.Details = "Payment processing"

	if payment.Amount <= 0 {
		transaction.ActivityType = models.FailedPayment
		transaction.Details = "Payment amount must be positive"
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("payment amount must be positive")
	}

	user, err := p.userRepo.FindByID(payment.CustomerID)
	if err != nil {
		transaction.ActivityType = models.FailedPayment
		transaction.Details = "Invalid customer ID"
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("invalid customer ID")
	}

	if !user.IsActive {
		transaction.ActivityType = models.FailedPayment
		transaction.Details = "Customer is not active"
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("customer is not active")
	}

	if user.Balance < payment.Amount {
		transaction.ActivityType = models.FailedPayment
		transaction.Details = "Insufficient balance"
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("insufficient balance")
	}

	user.Balance -= payment.Amount

	userRoles, err := p.roleRepo.FindRoleByUserID(user.ID)
	if err != nil {
		transaction.ActivityType = models.FailedPayment
		transaction.Details = "Error finding user roles"
		p.transactionRepo.CreateTransaction(transaction)
		return nil, err
	}

	for _, role := range *userRoles {
		if role.RoleID == "1" {
			user.Balance += payment.Amount
		}
	}

	if err := p.userRepo.UpdateUser(*user); err != nil {
		transaction.ActivityType = models.FailedPayment
		transaction.Details = "Failed to update customer"
		p.transactionRepo.CreateTransaction(transaction)
		return nil, err
	}
	transaction.ActivityType = models.PaymentActivity
	transaction.Details = "Payment processed successfully"

	trx, err := p.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}
	paymentResponse := mapper.TransactionModelToPaymentResponse(trx)

	if userRoles != nil {
		for _, role := range *userRoles {
			if role.RoleID == "1" {
				paymentResponse.MerchantID = payment.MerchantID
			}
			if role.RoleID == "2" {
				paymentResponse.CustomerID = payment.CustomerID
			}
		}
	}
	paymentResponse.Amount = payment.Amount

	return &paymentResponse, nil
}

func (p *transactionService) TransactionHistoryByUserID(userID string) ([]response.UserTransactionHistoryResponse, error) {
	transactions, err := p.transactionRepo.FindAllTransaction()
	if err != nil {
		return nil, err
	}

	user, err := p.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	var userTransactions []*models.Transaction
	for _, trx := range transactions {
		if trx.CustomerID == userID {
			userTransactions = append(userTransactions, &trx)
		}
	}

	userTransactionHistory := response.UserTransactionHistoryResponse{
		User:         *user,
		Transactions: userTransactions,
		TotalCount:   len(userTransactions),
	}

	return []response.UserTransactionHistoryResponse{userTransactionHistory}, nil
}
