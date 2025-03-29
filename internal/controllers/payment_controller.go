package controllers

import (
	"encoding/json"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/services"
	"net/http"

	"github.com/gorilla/mux"
)

type TransactionController struct {
	paymentService services.PaymentService
}

func NewTransactionController(paymentService services.PaymentService) TransactionController {
	return TransactionController{paymentService: paymentService}
}

func (t *TransactionController) Payment(w http.ResponseWriter, r *http.Request) {
	var request request.PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	payment, err := t.paymentService.ProcessPayment(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiRes := response.ApiResponse{
		Status:  http.StatusOK,
		Message: "Payment successful",
		Data:    payment,
	}
	response.CommonResponse(w, apiRes)
}

func (t *TransactionController) TransactionHistory(w http.ResponseWriter, r *http.Request) {
	customerID := mux.Vars(r)["id"]
	if customerID == "" {
		http.Error(w, "Customer ID is required", http.StatusBadRequest)
		return
	}
	transactions, err := t.paymentService.TransactionHistoryByUserID(customerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiRes := response.ApiResponse{
		Status:  http.StatusOK,
		Message: "Transaction history retrieved",
		Data:    transactions,
	}
	response.CommonResponse(w, apiRes)
}
