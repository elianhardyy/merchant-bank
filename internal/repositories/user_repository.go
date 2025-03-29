package repositories

import (
	"errors"
	"go-json/constant"
	"go-json/internal/models"
	"go-json/utils"
	"strconv"
	"sync"
)

type UserRepository interface {
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id string) (*models.User, error)
	CreateUser(User models.User, roleID []string) (*models.User, error)
	UpdateUser(User models.User) error
	FindAll() ([]models.User, error)
}

type userRepository struct {
	Users     []models.User
	roles     []models.Role
	userRoles []models.UserRole
	mu        sync.RWMutex
}

func NewUserRepository(users []models.User, roles []models.Role, userRoles []models.UserRole) UserRepository {
	return &userRepository{
		Users:     users,
		roles:     roles,
		userRoles: userRoles,
		mu:        sync.RWMutex{},
	}
}

func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	err := utils.ReadJSONFile(constant.USER_FILE, &r.Users)
	if err != nil {
		return nil, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, User := range r.Users {
		if User.Username == username {
			UserCopy := User
			return &UserCopy, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, User := range r.Users {
		if User.Email == email {
			UserCopy := User
			return &UserCopy, nil
		}
	}
	return nil, errors.New("user not found")

}

func (r *userRepository) FindByID(id string) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, User := range r.Users {
		if User.ID == id {
			UserCopy := User
			return &UserCopy, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepository) CreateUser(User models.User, roleIDs []string) (*models.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, existing := range r.Users {
		if existing.Username == User.Username {
			return nil, errors.New("username already exists")
		}
		if existing.Email == User.Email {
			return nil, errors.New("email already exists")
		}
	}
	newID := strconv.Itoa(len(r.Users) + 1)
	User.ID = newID
	User.Balance = 1000000.0
	User.IsActive = false

	var userRoles []models.UserRole
	roleFound := false

	for _, roleID := range roleIDs {
		for _, role := range r.roles {
			if role.ID == roleID {
				roleFound = true
				userRoles = append(userRoles, models.UserRole{
					ID:     strconv.Itoa(len(r.userRoles) + 1),
					UserID: User.ID,
					RoleID: roleID,
				})
			}
		}
	}

	if !roleFound {
		return nil, errors.New("invalid role ID")
	}
	r.Users = append(r.Users, User)
	r.userRoles = append(r.userRoles, userRoles...)

	if err := utils.WriteJSONFile(constant.USER_FILE, r.Users); err != nil {
		return nil, err
	}

	if err := utils.WriteJSONFile(constant.USER_ROLE_FILE, r.userRoles); err != nil {
		return nil, err
	}

	return &User, nil
}

func (r *userRepository) UpdateUser(User models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	found := false
	for i, existing := range r.Users {
		if existing.ID == User.ID {
			r.Users[i] = User
			found = true
			break
		}
	}

	if !found {
		return errors.New("user not found")
	}

	if err := utils.WriteJSONFile(constant.USER_FILE, r.Users); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.Users, nil
}
