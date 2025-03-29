package request

type PaymentRequest struct {
	CustomerID string  `json:"customer_id" validate:"required"`
	MerchantID string  `json:"merchant_id" validate:"required"`
	Amount     float64 `json:"amount" validate:"required"`
}
