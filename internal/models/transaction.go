package models

import "time"

type ActivityType string

const (
	LoginActivity   ActivityType = "LOGIN"
	PaymentActivity ActivityType = "PAYMENT"
	LogoutActivity  ActivityType = "LOGOUT"
	FailedLogin     ActivityType = "FAILED_LOGIN"
	FailedPayment   ActivityType = "FAILED_PAYMENT"
)

type Transaction struct {
	ID           string       `json:"id"`
	CustomerID   string       `json:"customer_id"`
	ActivityType ActivityType `json:"activity_type"`
	Timestamp    time.Time    `json:"timestamp"`
	Details      string       `json:"details"`
	Amount       float64      `json:"amount,omitempty"`
	MerchantID   string       `json:"merchant_id,omitempty"`
}
