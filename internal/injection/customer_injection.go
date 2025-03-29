package injection

import (
	"go-json/internal/controllers"
	"go-json/internal/models"
	"go-json/internal/repositories"
	"go-json/internal/security"
	"go-json/internal/services"
)

func InitCustomerAPI(customer []models.Customer, role []models.Role, userRole []models.UserRole, token security.TokenService) controllers.CustomerController {
	userRepository := repositories.NewCustomerRepository(customer, role, userRole)
	roleRepository := repositories.NewRoleRepository(role, userRole)
	userService := services.NewCustomerService(userRepository, roleRepository, token)
	return controllers.NewCustomerController(userService)
}
