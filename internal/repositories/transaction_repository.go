package repositories

import (
	"go-json/constant"
	"go-json/internal/models"
	"go-json/utils"
	"strconv"
	"sync"
)

type TransactionRepository interface {
	CreateTransaction(transaction models.Transaction) (*models.Transaction, error)
	FindAllTransaction() ([]models.Transaction, error)
}

type transactionRepository struct {
	transactions []models.Transaction
	mu           sync.RWMutex
}

func NewTransactionRepository(transactions []models.Transaction) TransactionRepository {
	return &transactionRepository{
		transactions: transactions,
		mu:           sync.RWMutex{},
	}
}

// CreateTransaction implements TransactionRepository.
func (t *transactionRepository) CreateTransaction(transaction models.Transaction) (*models.Transaction, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	newID := strconv.Itoa(len(t.transactions) + 1)
	transaction.ID = newID
	t.transactions = append(t.transactions, transaction)
	if err := utils.WriteJSONFile(constant.TRANSACTION_FILE, t.transactions); err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (t *transactionRepository) FindAllTransaction() ([]models.Transaction, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.transactions, nil
}
