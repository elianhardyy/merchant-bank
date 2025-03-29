package services

import (
	"errors"
	"go-json/internal/dtos/mapper"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/models"
	"go-json/internal/repositories"
	"time"
)

type PaymentService interface {
	ProcessPayment(payment request.PaymentRequest) (*response.PaymentResponse, error)
	TransactionHistoryByUserID(userID string) ([]response.UserTransactionHistoryResponse, error)
}

type paymentService struct {
	userRepo        repositories.CustomerRepository
	transactionRepo repositories.TransactionRepository
	roleRepo        repositories.RoleRepository
}

func NewPaymentService(userRepo repositories.CustomerRepository, transactionRepo repositories.TransactionRepository, roleRepo repositories.RoleRepository) PaymentService {
	return &paymentService{userRepo: userRepo, transactionRepo: transactionRepo, roleRepo: roleRepo}
}

func (p *paymentService) ProcessPayment(payment request.PaymentRequest) (*response.PaymentResponse, error) {
	var transaction models.Transaction
	transaction.CustomerID = payment.CustomerID
	transaction.MerchantID = payment.MerchantID
	transaction.Timestamp = time.Now()

	if payment.Amount <= 0 {
		transaction.ActivityType = models.FailedPayment
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("payment amount must be positive")
	}

	user, err := p.userRepo.FindByID(payment.CustomerID)
	if err != nil {
		transaction.ActivityType = models.FailedPayment
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("invalid customer ID")
	}

	if !user.IsActive {
		transaction.ActivityType = models.FailedPayment
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("customer is not active")
	}

	if user.Balance < payment.Amount {
		transaction.ActivityType = models.FailedPayment
		p.transactionRepo.CreateTransaction(transaction)
		return nil, errors.New("insufficient balance")
	}

	user.Balance -= payment.Amount

	userRoles, err := p.roleRepo.FindRoleByUserID(user.ID)
	if err != nil {
		return nil, err
	}

	for _, role := range *userRoles {
		if role.RoleID == "1" {
			user.Balance += payment.Amount
		}
	}

	if err := p.userRepo.UpdateCustomer(*user); err != nil {
		transaction.ActivityType = models.FailedPayment
		p.transactionRepo.CreateTransaction(transaction)
		return nil, err
	}

	transaction.ActivityType = models.PaymentActivity
	trx, err := p.transactionRepo.CreateTransaction(transaction)
	if err != nil {
		return nil, err
	}

	paymentResponse := response.PaymentResponse{
		CustomerID: payment.CustomerID,
		MerchantID: payment.MerchantID,
	}

	for _, role := range *userRoles {
		if role.RoleID == "1" {
			paymentResponse.MerchantID = payment.MerchantID
		}
		if role.RoleID == "2" {
			paymentResponse.CustomerID = payment.CustomerID
		}
	}

	trxResp := mapper.TransactionModelToPaymentResponse(trx)
	trxResp = paymentResponse
	return &trxResp, nil
}

func (p *paymentService) TransactionHistoryByUserID(userID string) ([]response.UserTransactionHistoryResponse, error) {
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
