package services

import (
	"errors"
	"go-json/internal/dtos/mapper"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"

	"go-json/internal/repositories"
	"go-json/internal/security"

	"github.com/go-playground/validator/v10"
)

type UserService interface {
	CreateUser(customer request.RegisterRequest) (*response.RegisterResponse, error)
	Login(login request.LoginRequest) (*response.LoginResponse, error)
	Logout(token string) (*response.LogoutResponse, error)
	Refresh(token string) (*response.RefreshResponse, error)
	FindAllUser() ([]*response.UserResponse, error)
	//ProfileCustomer(id string)(*response.ProfileCustomerResponse, error)
}

type userService struct {
	userRepo repositories.UserRepository
	roleRepo repositories.RoleRepository
	token    security.TokenService
}

func NewUserService(user repositories.UserRepository, role repositories.RoleRepository, token security.TokenService) UserService {
	return &userService{
		userRepo: user,
		roleRepo: role,
		token:    token,
	}
}

func (s *userService) CreateUser(user request.RegisterRequest) (*response.RegisterResponse, error) {
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		return nil, err
	}
	userModel := mapper.UserRequestToModel(user)

	var roleIDs []string
	for _, roleName := range user.Role {
		role, err := s.roleRepo.FindByRoleName(roleName)
		if err != nil {
			return nil, err
		}
		roleIDs = append(roleIDs, role.ID)
	}

	users, err := s.userRepo.CreateUser(userModel, roleIDs)
	if err != nil {
		return nil, err
	}
	userResponse := mapper.UserModelToResponse(*users)
	return &userResponse, nil
}

func (s *userService) Login(login request.LoginRequest) (*response.LoginResponse, error) {
	validate := validator.New()
	err := validate.Struct(login)
	if err != nil {
		return nil, err
	}
	loginReq := mapper.LoginRequestToModel(login)
	customer, err := s.userRepo.FindByUsername(loginReq.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if customer.Password != login.Password {
		return nil, errors.New("invalid credentials")
	}

	customer.IsActive = true
	err = s.userRepo.UpdateUser(*customer)
	if err != nil {
		return nil, err
	}

	userRoles, err := s.roleRepo.FindRoleByUserID(customer.ID)
	if err != nil {
		return nil, err
	}

	var roleNames []string
	for _, userRole := range *userRoles {
		role, err := s.roleRepo.FindByRoleID(userRole.RoleID)
		if err != nil {
			return nil, err
		}
		roleNames = append(roleNames, role.Name)
	}

	token, err := s.token.GenerateToken(customer.Username, customer.Email, roleNames...)
	if err != nil {
		return nil, err
	}

	customerResponse := mapper.UserModelToLoginResponse(token)
	return &customerResponse, nil
}

func (s *userService) Logout(token string) (*response.LogoutResponse, error) {
	logoutResponse := mapper.ToLogoutResponse(token)
	claims, err := s.token.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	user, err := s.userRepo.FindByEmail(claims.Email)
	if err != nil {
		return nil, err
	}
	user.IsActive = false
	err = s.userRepo.UpdateUser(*user)
	if err != nil {
		return nil, err
	}
	return &logoutResponse, nil
}

func (s *userService) Refresh(token string) (*response.RefreshResponse, error) {
	refreshResponse := mapper.ToRefreshResponse(token)
	return &refreshResponse, nil
}

func (s *userService) FindAllUser() ([]*response.UserResponse, error) {
	customers, err := s.userRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var userResponses []*response.UserResponse
	for _, customer := range customers {
		userRoles, err := s.roleRepo.FindRoleByUserID(customer.ID)
		if err != nil {
			return nil, err
		}

		var roleNames []string
		for _, userRole := range *userRoles {
			role, err := s.roleRepo.FindByRoleID(userRole.RoleID)
			if err != nil {
				return nil, err
			}
			roleNames = append(roleNames, role.Name)
		}

		userResponse := mapper.UserModelToUserResponse(customer, roleNames)
		userResponses = append(userResponses, &userResponse)
	}
	return userResponses, nil
}
