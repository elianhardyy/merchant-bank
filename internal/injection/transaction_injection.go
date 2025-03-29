package injection

import (
	"go-json/internal/controllers"
	"go-json/internal/models"
	"go-json/internal/repositories"
	"go-json/internal/services"
)

func InitTransactionAPI(transaction []models.Transaction, User []models.User, role []models.Role, userRole []models.UserRole) controllers.TransactionController {
	transactionRepository := repositories.NewTransactionRepository(transaction)
	roleRepository := repositories.NewRoleRepository(role, userRole)
	userRepository := repositories.NewUserRepository(User, role, userRole)
	transactionService := services.NewTransactionService(userRepository, transactionRepository, roleRepository)
	return controllers.NewTransactionController(transactionService)
}
