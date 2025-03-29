package injection

import (
	"go-json/internal/controllers"
	"go-json/internal/models"
	"go-json/internal/repositories"
	"go-json/internal/services"
)

func InitTransactionAPI(transaction []models.Transaction, customer []models.Customer, role []models.Role, userRole []models.UserRole) controllers.TransactionController {
	transactionRepository := repositories.NewTransactionRepository(transaction)
	roleRepository := repositories.NewRoleRepository(role, userRole)
	userRepository := repositories.NewCustomerRepository(customer, role, userRole)
	transactionService := services.NewPaymentService(userRepository, transactionRepository, roleRepository)
	return controllers.NewTransactionController(transactionService)
}
