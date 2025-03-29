package injection

import (
	"go-json/internal/controllers"
	"go-json/internal/models"
	"go-json/internal/repositories"
	"go-json/internal/security"
	"go-json/internal/services"
)

func InitUserAPI(User []models.User, role []models.Role, userRole []models.UserRole, token security.TokenService) controllers.UserController {
	userRepository := repositories.NewUserRepository(User, role, userRole)
	roleRepository := repositories.NewRoleRepository(role, userRole)
	userService := services.NewUserService(userRepository, roleRepository, token)
	return controllers.NewUserController(userService)
}
