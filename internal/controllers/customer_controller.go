package controllers

import (
	"encoding/json"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/services"
	"net/http"
)

type CustomerController struct {
	customerService services.CustomerService
}

func NewCustomerController(customerService services.CustomerService) CustomerController {
	return CustomerController{customerService: customerService}
}

func (c *CustomerController) Register(w http.ResponseWriter, r *http.Request) {
	var request request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	customer, err := c.customerService.CreateCustomer(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiRes := response.ApiResponse{
		Status:  http.StatusCreated,
		Message: "Customer registered successfully",
		Data:    customer,
	}
	response.CommonResponse(w, apiRes)
}

func (c *CustomerController) Login(w http.ResponseWriter, r *http.Request) {
	var request request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := c.customerService.Login(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiRes := response.ApiResponse{
		Status:  http.StatusOK,
		Message: "Login successful",
		Data:    token,
	}
	response.CommonResponse(w, apiRes)
}

func (c *CustomerController) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	logoutResponse, err := c.customerService.Logout(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiRes := response.ApiResponse{
		Status:  http.StatusOK,
		Message: logoutResponse.Message,
		Data:    logoutResponse,
	}
	response.CommonResponse(w, apiRes)
}
