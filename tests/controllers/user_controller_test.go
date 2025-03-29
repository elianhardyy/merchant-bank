package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-json/internal/controllers"
	"go-json/internal/dtos/request"
	"go-json/internal/dtos/response"
	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(user request.RegisterRequest) (*response.RegisterResponse, error) {
	args := m.Called(user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.RegisterResponse), args.Error(1)
}

func (m *MockUserService) Login(login request.LoginRequest) (*response.LoginResponse, error) {
	args := m.Called(login)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.LoginResponse), args.Error(1)
}

func (m *MockUserService) Logout(token string) (*response.LogoutResponse, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.LogoutResponse), args.Error(1)
}

func (m *MockUserService) Refresh(token string) (*response.RefreshResponse, error) {
	args := m.Called(token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*response.RefreshResponse), args.Error(1)
}

func (m *MockUserService) FindAllUser() ([]*response.UserResponse, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*response.UserResponse), args.Error(1)
}

type UserControllerTestSuite struct {
	suite.Suite
	userService *MockUserService
	controller  controllers.UserController
}

func (suite *UserControllerTestSuite) SetupTest() {
	suite.userService = new(MockUserService)
	suite.controller = controllers.NewUserController(suite.userService)
}

func (suite *UserControllerTestSuite) TestRegister() {
	registerReq := request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     []string{"customer"},
	}

	registerResp := &response.RegisterResponse{
		ID:       "1",
		Username: "testuser",
		Email:    "test@example.com",
	}

	suite.userService.On("CreateUser", mock.MatchedBy(func(req request.RegisterRequest) bool {
		return req.Username == registerReq.Username && req.Email == registerReq.Email
	})).Return(registerResp, nil)

	reqBody, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Register)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var apiResp response.ApiResponse
	err := json.Unmarshal(rr.Body.Bytes(), &apiResp)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusCreated, apiResp.Status)
	assert.Equal(suite.T(), "Customer registered successfully", apiResp.Message)

	suite.userService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestRegisterError() {
	registerReq := request.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
		Role:     []string{"customer"},
	}

	suite.userService.On("CreateUser", mock.MatchedBy(func(req request.RegisterRequest) bool {
		return req.Username == registerReq.Username && req.Email == registerReq.Email
	})).Return(nil, errors.New("registration error"))

	reqBody, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Register)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)

	suite.userService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestLogin() {
	loginReq := request.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}

	loginResp := &response.LoginResponse{
		Token: "test-token",
	}

	suite.userService.On("Login", mock.MatchedBy(func(req request.LoginRequest) bool {
		return req.Username == loginReq.Username && req.Password == loginReq.Password
	})).Return(loginResp, nil)

	reqBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var apiResp response.ApiResponse
	err := json.Unmarshal(rr.Body.Bytes(), &apiResp)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, apiResp.Status)
	assert.Equal(suite.T(), "Login successful", apiResp.Message)

	suite.userService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestLoginError() {
	loginReq := request.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword",
	}

	suite.userService.On("Login", mock.MatchedBy(func(req request.LoginRequest) bool {
		return req.Username == loginReq.Username && req.Password == loginReq.Password
	})).Return(nil, errors.New("invalid credentials"))

	reqBody, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Login)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, rr.Code)

	suite.userService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestLogout() {
	token := "test-token"
	logoutResp := &response.LogoutResponse{
		Message: "Logout successful",
	}

	suite.userService.On("Logout", token).Return(logoutResp, nil)

	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Authorization", token)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Logout)
	handler.ServeHTTP(rr, req)

	assert.Equal(suite.T(), http.StatusOK, rr.Code)

	var apiResp response.ApiResponse
	err := json.Unmarshal(rr.Body.Bytes(), &apiResp)
	assert.NoError(suite.T(), err)

	assert.Equal(suite.T(), http.StatusOK, apiResp.Status)
	assert.Equal(suite.T(), "Logout successful", apiResp.Message)

	suite.userService.AssertExpectations(suite.T())
}

func (suite *UserControllerTestSuite) TestLogoutNoToken() {
	req, _ := http.NewRequest("POST", "/logout", nil)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(suite.controller.Logout)
	handler.ServeHTTP(rr, req)
}
