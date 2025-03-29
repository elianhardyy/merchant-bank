package repositories

import (
	"errors"
	"go-json/constant"
	"go-json/internal/models"
	"go-json/utils"
	"strconv"
	"sync"
)

type CustomerRepository interface {
	FindByUsername(username string) (*models.Customer, error)
	FindByEmail(email string) (*models.Customer, error)
	FindByID(id string) (*models.Customer, error)
	CreateCustomer(customer models.Customer, roleID []string) (*models.Customer, error)
	UpdateCustomer(customer models.Customer) error
	FindAll() ([]models.Customer, error)
}

type customerRepository struct {
	customers []models.Customer
	roles     []models.Role
	userRoles []models.UserRole
	mu        sync.RWMutex
}

func NewCustomerRepository(customers []models.Customer, roles []models.Role, userRoles []models.UserRole) CustomerRepository {
	return &customerRepository{
		customers: customers,
		roles:     roles,
		userRoles: userRoles,
		mu:        sync.RWMutex{},
	}
}

func (r *customerRepository) FindByUsername(username string) (*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, customer := range r.customers {
		if customer.Username == username {
			customerCopy := customer
			return &customerCopy, nil
		}
	}
	return nil, errors.New("customer not found")
}

func (r *customerRepository) FindByEmail(email string) (*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, customer := range r.customers {
		if customer.Email == email {
			customerCopy := customer
			return &customerCopy, nil
		}
	}
	return nil, errors.New("customer not found")

}

func (r *customerRepository) FindByID(id string) (*models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, customer := range r.customers {
		if customer.ID == id {
			customerCopy := customer
			return &customerCopy, nil
		}
	}
	return nil, errors.New("customer not found")
}

func (r *customerRepository) CreateCustomer(customer models.Customer, roleIDs []string) (*models.Customer, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.customers {
		if existing.Username == customer.Username {
			return nil, errors.New("username already exists")
		}
		if existing.Email == customer.Email {
			return nil, errors.New("email already exists")
		}
	}
	newID := strconv.Itoa(len(r.customers) + 1)
	customer.ID = newID
	customer.Balance = 0.0
	customer.IsActive = false

	var userRoles []models.UserRole
	roleFound := false

	for _, roleID := range roleIDs {
		for _, role := range r.roles {
			if role.ID == roleID {
				roleFound = true
				userRoles = append(userRoles, models.UserRole{
					ID:     strconv.Itoa(len(r.userRoles) + 1),
					UserID: customer.ID,
					RoleID: roleID,
				})
			}
		}
	}

	if !roleFound {
		return nil, errors.New("invalid role ID")
	}

	r.customers = append(r.customers, customer)
	r.userRoles = append(r.userRoles, userRoles...)

	if err := utils.WriteJSONFile(constant.CUSTOMER_FILE, r.customers); err != nil {
		return nil, err
	}

	if err := utils.WriteJSONFile(constant.USER_ROLE_FILE, r.userRoles); err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *customerRepository) UpdateCustomer(customer models.Customer) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	found := false
	for i, existing := range r.customers {
		if existing.ID == customer.ID {
			r.customers[i] = customer
			found = true
			break
		}
	}

	if !found {
		return errors.New("user not found")
	}

	if err := utils.WriteJSONFile(constant.CUSTOMER_FILE, r.customers); err != nil {
		return err
	}

	return nil
}

func (r *customerRepository) FindAll() ([]models.Customer, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.customers, nil
}
