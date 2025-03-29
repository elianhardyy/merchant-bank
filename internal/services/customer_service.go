package services

import (
	"errors"
	"go-json/internal/dtos/mapper"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"

	"go-json/internal/repositories"
	"go-json/internal/security"
	"go-json/internal/validators"
)

type CustomerService interface {
	CreateCustomer(customer request.RegisterRequest) (*response.RegisterResponse, error)
	Login(login request.LoginRequest) (*response.LoginResponse, error)
	Logout(token string) (*response.LogoutResponse, error)
	Refresh(token string) (*response.RefreshResponse, error)
	//ProfileCustomer(id string)(*response.ProfileCustomerResponse, error)
}

type customerService struct {
	customerRepo repositories.CustomerRepository
	roleRepo     repositories.RoleRepository
	validator    *validators.CustomerValidator
	token        security.TokenService
}

func NewCustomerService(repo repositories.CustomerRepository, role repositories.RoleRepository, token security.TokenService) CustomerService {
	return &customerService{
		customerRepo: repo,
		roleRepo:     role,
		validator:    validators.NewCustomerValidator(),
		token:        token,
	}
}

func (s *customerService) CreateCustomer(customer request.RegisterRequest) (*response.RegisterResponse, error) {
	if err := s.validator.Validate(customer); err != nil {
		return nil, err
	}
	customerModel := mapper.CustomerRequestToModel(customer)

	var roleIDs []string
	for _, roleName := range customer.Role {
		role, err := s.roleRepo.FindByRoleName(roleName)
		if err != nil {
			return nil, err
		}
		roleIDs = append(roleIDs, role.ID)
	}

	customers, err := s.customerRepo.CreateCustomer(customerModel, roleIDs)
	if err != nil {
		return nil, err
	}
	customerResponse := mapper.CustomerModelToResponse(*customers)
	return &customerResponse, nil
}

func (s *customerService) Login(login request.LoginRequest) (*response.LoginResponse, error) {
	loginReq := mapper.LoginRequestToModel(login)
	customer, err := s.customerRepo.FindByUsername(loginReq.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if customer.Password != login.Password {
		return nil, errors.New("invalid credentials")
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
	customer.IsActive = true
	customerResponse := mapper.CustomerModelToLoginResponse(token)
	return &customerResponse, nil
}

func (s *customerService) Logout(token string) (*response.LogoutResponse, error) {
	logoutResponse := mapper.ToLogoutResponse(token)
	claims, err := s.token.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	user, err := s.customerRepo.FindByEmail(claims.Email)
	if err != nil {
		return nil, err
	}
	user.IsActive = false
	return &logoutResponse, nil
}

func (s *customerService) Refresh(token string) (*response.RefreshResponse, error) {
	refreshResponse := mapper.ToRefreshResponse(token)
	return &refreshResponse, nil
}
