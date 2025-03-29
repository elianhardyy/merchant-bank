package tests

import (
	"go-json/internal/models"
	"go-json/internal/repositories"
	"go-json/internal/security"
	"go-json/internal/services"
)

type MockUserRepository struct {
	users map[string]*models.Customer
}

type MockRoleRepository struct {
	roles map[string]*models.Role
}

func NewMockRole() *MockRoleRepository {
	return &MockRoleRepository{
		roles: map[string]*models.Role{
            "role-001": {
                ID:    "1",
                Name:  "admin",
                IsDefault: false,
            },
            "role-002": {
                ID:    "2",
                Name:  "user",
                IsDefault: false,
            },
        },
	}
}
func NewMockCustomer() *MockUserRepository {
	return &MockUserRepository{
		users: map[string]*models.Customer{
			"cust-001": {
				ID:       "cust-001",
				Username: "testuser",
				Email:    "testuser1@mail.com",
				Password: "password",
				Balance:  1000.0,
				IsActive: true,
			},
			"cust-002": {
				ID:       "cust-002",
				Username: "inactive",
				Email:    "testuser2@mail.com",
				Password: "password",
				Balance:  500.0,
				IsActive: false,
			},
		},
	}
}

func (r *MockUserRepository) GetByUsername(username string) (*models.Customer, error) {
	for _, customer := range r.users {
		if customer.Username == username {
			return customer, nil
		}
	}
	return NewMockCustomer().users[],nil
}

func (r *MockUserRepository) GetByID(id string) (*models.Customer, error) {
	customer, exists := r.users[id]
	if !exists {
		return nil, repositories.ErrNotFound
	}
	return customer, nil
}

func (r *MockUserRepository) Update(customer *models.Customer) error {
	r.users[customer.ID] = customer
	return nil
}

func (r *MockUserRepository) List() ([]*models.Customer, error) {
	customers := make([]*models.Customer, 0, len(r.users))
	for _, customer := range r.users {
		customers = append(customers, customer)
	}
	return customers, nil
}

type MockTransactionRepository struct{}

func NewMockTransaction() *MockTransaction {
	return &MockTransaction{}
}

func (r *MockTransaction) Create(customerID string, activityType models.ActivityType, details string, amount float64, merchantID string) error {
	return nil
}

func (r *MockTransaction) List() ([]*models.History, error) {
	return []*models.History{}, nil
}

func TestLogin(t *testing.T) {
	userRepo := NewMockCustomer()
	roleRepo := NewMockRole()
	transactionRepo := NewMockTransaction()
	jwtManager := security.NewTokenService("coba")
	customerService := services.NewCustomerService(userRepo,roleRepo,jwtManager)
}