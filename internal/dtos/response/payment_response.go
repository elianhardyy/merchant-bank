package response

import (
	"go-json/internal/models"
	"time"
)

type PaymentResponse struct {
	ID           string    `json:"id"`
	CustomerID   string    `json:"customer_id"`
	ActivityType string    `json:"activity_type"`
	Timestamp    time.Time `json:"timestamp"`
	Details      string    `json:"details"`
	Amount       float64   `json:"amount"`
	MerchantID   string    `json:"merchant_id"`
}

type UserTransactionHistoryResponse struct {
	User         models.User           `json:"user,omitempty"`
	Transactions []*models.Transaction `json:"transactions,omitempty"`
	TotalCount   int                   `json:"total_count"`
}
