package services_test

import (
	"errors"
	"go-json/internal/dtos/request"
	"go-json/internal/models"
	"go-json/internal/security"
	"go-json/internal/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user models.User, roleID []string) (*models.User, error) {
	args := m.Called(user, roleID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateUser(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll() ([]models.User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.User), args.Error(1)
}

type MockRoleRepository struct {
	mock.Mock
}

func (m *MockRoleRepository) FindByRoleName(name string) (*models.Role, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) FindByRoleID(id string) (*models.Role, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *MockRoleRepository) FindRoleByUserID(userID string) (*[]models.UserRole, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]models.UserRole), args.Error(1)
}

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(username, email string, role ...string) (string, error) {
	args := m.Called(username, email, role)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) IsTokenExpired(tokenString string) bool {
	args := m.Called(tokenString)
	return args.Bool(0)
}

func (m *MockTokenService) IsTokenValid(tokenString string) bool {
	args := m.Called(tokenString)
	return args.Bool(0)
}

func (m *MockTokenService) VerifyToken(tokenString string) (*security.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*security.Claims), args.Error(1)
}

type UserServiceTestSuite struct {
	suite.Suite
	userRepo  *MockUserRepository
	roleRepo  *MockRoleRepository
	tokenSvc  *MockTokenService
	userSvc   services.UserService
	testUser  models.User
	testRoles []models.Role
}

func (suite *UserServiceTestSuite) SetupTest() {
	suite.userRepo = new(MockUserRepository)
	suite.roleRepo = new(MockRoleRepository)
	suite.tokenSvc = new(MockTokenService)
	suite.userSvc = services.NewUserService(suite.userRepo, suite.roleRepo, suite.tokenSvc)

	suite.testUser = models.User{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Balance:  1000.0,
		IsActive: false,
	}

	suite.testRoles = []models.Role{
		{
			ID:        "1",
			Name:      "merchant",
			IsDefault: false,
		},
		{
			ID:        "2",
			Name:      "customer",
			IsDefault: true,
		},
	}
}

func (suite *UserServiceTestSuite) TestCreateUser() {
	registerReq := request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     []string{"customer"},
	}

	suite.roleRepo.On("FindByRoleName", "customer").Return(&suite.testRoles[1], nil)

	suite.userRepo.On("CreateUser", mock.AnythingOfType("models.User"), []string{"2"}).Return(&suite.testUser, nil)

	response, err := suite.userSvc.CreateUser(registerReq)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), suite.testUser.ID, response.ID)
	assert.Equal(suite.T(), suite.testUser.Username, response.Username)
	assert.Equal(suite.T(), suite.testUser.Email, response.Email)

	suite.roleRepo.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestCreateUserRoleError() {
	registerReq := request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     []string{"invalid_role"},
	}

	suite.roleRepo.On("FindByRoleName", "invalid_role").Return(nil, errors.New("role not found"))

	response, err := suite.userSvc.CreateUser(registerReq)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), "role not found", err.Error())

	suite.roleRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLogin() {
	loginReq := request.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	userRoles := []models.UserRole{
		{
			ID:     "1",
			UserID: "1",
			RoleID: "2",
		},
	}

	suite.userRepo.On("FindByUsername", "testuser").Return(&suite.testUser, nil)
	suite.userRepo.On("UpdateUser", mock.AnythingOfType("models.User")).Return(nil)
	suite.roleRepo.On("FindRoleByUserID", "1").Return(&userRoles, nil)
	suite.roleRepo.On("FindByRoleID", "2").Return(&suite.testRoles[1], nil)
	suite.tokenSvc.On("GenerateToken", "testuser", "test@example.com", mock.AnythingOfType("[]string")).Return("test-token", nil)

	response, err := suite.userSvc.Login(loginReq)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "test-token", response.Token)

	suite.userRepo.AssertExpectations(suite.T())
	suite.roleRepo.AssertExpectations(suite.T())
	suite.tokenSvc.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLoginInvalidCredentials() {
	loginReq := request.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	suite.userRepo.On("FindByUsername", "testuser").Return(&suite.testUser, nil)

	response, err := suite.userSvc.Login(loginReq)

	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), response)
	assert.Equal(suite.T(), "invalid credentials", err.Error())

	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestLogout() {
	token := "test-token"
	claims := &security.Claims{
		Email: "test@example.com",
		Role:  []string{"customer"},
	}

	suite.tokenSvc.On("VerifyToken", token).Return(claims, nil)
	suite.userRepo.On("FindByEmail", "test@example.com").Return(&suite.testUser, nil)
	suite.userRepo.On("UpdateUser", mock.AnythingOfType("models.User")).Return(nil)

	response, err := suite.userSvc.Logout(token)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), "Logout successful", response.Message)

	suite.tokenSvc.AssertExpectations(suite.T())
	suite.userRepo.AssertExpectations(suite.T())
}

func (suite *UserServiceTestSuite) TestRefresh() {
	token := "test-token"

	response, err := suite.userSvc.Refresh(token)

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Equal(suite.T(), token, response.Token)
}

func (suite *UserServiceTestSuite) TestFindAllUser() {
	users := []models.User{suite.testUser}
	userRoles := []models.UserRole{
		{
			ID:     "1",
			UserID: "1",
			RoleID: "2",
		},
	}

	suite.userRepo.On("FindAll").Return(users, nil)
	suite.roleRepo.On("FindRoleByUserID", "1").Return(&userRoles, nil)
	suite.roleRepo.On("FindByRoleID", "2").Return(&suite.testRoles[1], nil)

	response, err := suite.userSvc.FindAllUser()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), response)
	assert.Len(suite.T(), response, 1)
	assert.Equal(suite.T(), suite.testUser.Username, response[0].Username)
	// assert.Equal(suite.T(), []string{"customer"}, response[0])

	suite.userRepo.AssertExpectations(suite.T())
	suite.roleRepo.AssertExpectations(suite.T())
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
