package mapper

import (
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"go-json/internal/models"
)

func UserRequestToModel(request request.RegisterRequest) models.User {
	return models.User{
		Username: request.Username,
		Email:    request.Email,
		Password: request.Password,
	}
}

func UserModelToResponse(User models.User) response.RegisterResponse {
	return response.RegisterResponse{
		ID:       User.ID,
		Username: User.Username,
		Email:    User.Email,
		Balance:  User.Balance,
		IsActive: User.IsActive,
	}
}

func LoginRequestToModel(request request.LoginRequest) models.User {
	return models.User{
		Username: request.Username,
		Password: request.Password,
	}
}

func UserModelToLoginResponse(token string) response.LoginResponse {
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

func UserModelToUserResponse(user models.User, roles []string) response.UserResponse {
	return response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Roles:    roles,
		Balance:  user.Balance,
		IsActive: user.IsActive,
	}
}
