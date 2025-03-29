package request

type PaymentRequest struct {
	CustomerID string  `json:"customer_id"`
	MerchantID string  `json:"merchant_id"`
	Amount     float64 `json:"amount"`
}
