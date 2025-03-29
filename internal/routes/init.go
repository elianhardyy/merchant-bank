package routes

import (
	"go-json/constant"
	"go-json/internal/injection"
	"go-json/internal/models"
	"go-json/internal/security"
	"go-json/utils"
	"os"

	"github.com/gorilla/mux"
)

var R = mux.NewRouter()

func InitRoute() {
	secret := []byte(os.Getenv("JWT_SECRET"))
	token := security.NewTokenService(secret)
	var customers []models.Customer
	err := utils.ReadJSONFile(constant.CUSTOMER_FILE, &customers)
	if err != nil {
		customers = []models.Customer{}
	}

	var roles []models.Role
	err = utils.ReadJSONFile(constant.ROLE_FILE, &roles)
	if err != nil {
		roles = []models.Role{}
	}

	var userRoles []models.UserRole
	err = utils.ReadJSONFile(constant.USER_ROLE_FILE, &userRoles)
	if err != nil {
		userRoles = []models.UserRole{}
	}

	var transactions []models.Transaction
	err = utils.ReadJSONFile(constant.TRANSACTION_FILE, &transactions)
	if err != nil {
		transactions = []models.Transaction{}
	}

	customerApi := injection.InitCustomerAPI(customers, roles, userRoles, token)
	UserRoutes(customerApi, token)

	transactionApi := injection.InitTransactionAPI(transactions, customers, roles, userRoles)
	TransactionRoutes(transactionApi, token)
}
