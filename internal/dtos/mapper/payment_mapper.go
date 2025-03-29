package mapper

import (
	"go-json/internal/dtos/response"
	"go-json/internal/models"
)

func TransactionModelToPaymentResponse(trx *models.Transaction) response.PaymentResponse {
	return response.PaymentResponse{
		ID:           trx.ID,
		Amount:       trx.Amount,
		CustomerID:   trx.CustomerID,
		MerchantID:   trx.MerchantID,
		ActivityType: string(trx.ActivityType),
		Details:      trx.Details,
		Timestamp:    trx.Timestamp,
	}
}
