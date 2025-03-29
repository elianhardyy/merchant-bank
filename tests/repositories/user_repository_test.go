package repositories_test

import (
	"go-json/internal/models"
	"go-json/internal/repositories"
	constant_test "go-json/tests/constant"
	"go-json/utils"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepositoryTestSuite struct {
	suite.Suite
	repo      repositories.UserRepository
	users     []models.User
	roles     []models.Role
	userRoles []models.UserRole
	testUser  models.User
}

func (suite *UserRepositoryTestSuite) SetupTest() {
	suite.users = []models.User{
		{
			ID:       "1",
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
			Balance:  1000.0,
			IsActive: false,
		},
	}

	suite.roles = []models.Role{
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

	suite.userRoles = []models.UserRole{
		{
			ID:     "1",
			UserID: "1",
			RoleID: "2",
		},
	}

	suite.testUser = models.User{
		Username: "newuser",
		Email:    "new@example.com",
		Password: "password123",
	}

	err := utils.WriteJSONFile(constant_test.USER_FILE, suite.users)
	assert.NoError(suite.T(), err)

	err = utils.WriteJSONFile(constant_test.ROLE_FILE, suite.roles)
	assert.NoError(suite.T(), err)

	err = utils.WriteJSONFile(constant_test.USER_ROLE_FILE, suite.userRoles)
	assert.NoError(suite.T(), err)

	suite.repo = repositories.NewUserRepository(suite.users, suite.roles, suite.userRoles)
}

func (suite *UserRepositoryTestSuite) TearDownTest() {
	os.Remove(constant_test.USER_FILE)
	os.Remove(constant_test.ROLE_FILE)
	os.Remove(constant_test.USER_ROLE_FILE)
}

func (suite *UserRepositoryTestSuite) TestFindByUsername() {
	user, err := suite.repo.FindByUsername("testuser")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "testuser", user.Username)

	user, err = suite.repo.FindByUsername("nonexistent")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

func (suite *UserRepositoryTestSuite) TestFindByEmail() {
	user, err := suite.repo.FindByEmail("test@example.com")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "test@example.com", user.Email)

	user, err = suite.repo.FindByEmail("nonexistent@example.com")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

func (suite *UserRepositoryTestSuite) TestFindByID() {
	user, err := suite.repo.FindByID("1")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "1", user.ID)

	user, err = suite.repo.FindByID("999")
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "user not found", err.Error())
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
	user, err := suite.repo.CreateUser(suite.testUser, []string{"2"})
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "newuser", user.Username)
	assert.Equal(suite.T(), "new@example.com", user.Email)
	assert.Equal(suite.T(), "2", user.ID) // Since it's the second user in the array
	assert.Equal(suite.T(), 1000000.0, user.Balance)
	assert.False(suite.T(), user.IsActive)

	duplicateUser := suite.testUser
	duplicateUser.Email = "another@example.com"
	user, err = suite.repo.CreateUser(duplicateUser, []string{"2"})
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "username already exists", err.Error())

	duplicateUser = suite.testUser
	duplicateUser.Username = "anotheruser"
	user, err = suite.repo.CreateUser(duplicateUser, []string{"2"})
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "email already exists", err.Error())

	user, err = suite.repo.CreateUser(models.User{
		Username: "validuser",
		Email:    "valid@example.com",
		Password: "password123",
	}, []string{"999"})
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), user)
	assert.Equal(suite.T(), "invalid role ID", err.Error())
}

func (suite *UserRepositoryTestSuite) TestUpdateUser() {
	user, err := suite.repo.FindByID("1")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)

	user.Username = "updateduser"
	user.Balance = 2000.0
	user.IsActive = true

	err = suite.repo.UpdateUser(*user)
	assert.NoError(suite.T(), err)

	updatedUser, err := suite.repo.FindByID("1")
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updatedUser)
	assert.Equal(suite.T(), "updateduser", updatedUser.Username)
	assert.Equal(suite.T(), 2000.0, updatedUser.Balance)
	assert.True(suite.T(), updatedUser.IsActive)

	nonExistingUser := models.User{
		ID:       "999",
		Username: "nonexistent",
		Email:    "nonexistent@example.com",
	}
	err = suite.repo.UpdateUser(nonExistingUser)
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user not found", err.Error())
}

func (suite *UserRepositoryTestSuite) TestFindAll() {
	users, err := suite.repo.FindAll()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), users)
	assert.Len(suite.T(), users, 1)
	assert.Equal(suite.T(), "testuser", users[0].Username)

	_, err = suite.repo.CreateUser(suite.testUser, []string{"2"})
	assert.NoError(suite.T(), err)

	users, err = suite.repo.FindAll()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), users)
	assert.Len(suite.T(), users, 2)
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
