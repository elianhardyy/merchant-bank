package controllers

import (
	"encoding/json"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/services"
	"net/http"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService: userService}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var request request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user, err := c.userService.CreateUser(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiRes := response.ApiResponse{
		Status:  http.StatusCreated,
		Message: "Customer registered successfully",
		Data:    user,
	}
	response.CommonResponse(w, apiRes)
}

func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var request request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := c.userService.Login(request)
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

func (c *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Authorization header required", http.StatusUnauthorized)
		return
	}
	logoutResponse, err := c.userService.Logout(token)
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

func (c *UserController) UserList(w http.ResponseWriter, r *http.Request) {
	customers, err := c.userService.FindAllUser()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiRes := response.ApiResponse{
		Status:  http.StatusOK,
		Message: "User list retrieved",
		Data:    customers,
	}
	response.CommonResponse(w, apiRes)
}
