package mapper

import (
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/models"
)

func CustomerRequestToModel(request request.RegisterRequest) models.Customer {
	return models.Customer{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}

func CustomerModelToResponse(customer models.Customer) response.RegisterResponse {
	return response.RegisterResponse{
		ID:       customer.ID,
		Username: customer.Username,
		Email:    customer.Email,
		Balance:  customer.Balance,
		IsActive: customer.IsActive,
	}
}

func LoginRequestToModel(request request.LoginRequest) models.Customer {
	return models.Customer{
		Username: request.Username,
		Password: request.Password,
	}
}

func CustomerModelToLoginResponse(token string) response.LoginResponse {
	return response.LoginResponse{
		Token: token,
	}
}

func ToLogoutResponse(token string) response.LogoutResponse {
	return response.LogoutResponse{
		Message: "Successfully logged out",
		Token:   token,
	}
}

func ToRefreshResponse(token string) response.RefreshResponse {
	return response.RefreshResponse{
		Token: token,
	}
}
