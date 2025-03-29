package repositories_test

import (
	"go-json/internal/models"
	"go-json/internal/repositories"
	constant_test "go-json/tests/constant"
	"go-json/utils"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TransactionRepositoryTestSuite struct {
	suite.Suite
	repo         repositories.TransactionRepository
	transactions []models.Transaction
	testTrx      models.Transaction
}

func (suite *TransactionRepositoryTestSuite) SetupTest() {
	suite.transactions = []models.Transaction{
		{
			ID:           "1",
			CustomerID:   "1",
			ActivityType: models.PaymentActivity,
			Timestamp:    time.Now(),
			Details:      "Test transaction",
			Amount:       100.0,
			MerchantID:   "2",
		},
	}

	suite.testTrx = models.Transaction{
		CustomerID:   "2",
		ActivityType: models.PaymentActivity,
		Timestamp:    time.Now(),
		Details:      "New test transaction",
		Amount:       200.0,
		MerchantID:   "3",
	}

	err := utils.WriteJSONFile(constant_test.TRANSACTION_FILE, suite.transactions)
	assert.NoError(suite.T(), err)

	suite.repo = repositories.NewTransactionRepository(suite.transactions)
}

func (suite *TransactionRepositoryTestSuite) TearDownTest() {
	os.Remove(constant_test.TRANSACTION_FILE)
}

func (suite *TransactionRepositoryTestSuite) TestCreateTransaction() {
	transaction, err := suite.repo.CreateTransaction(suite.testTrx)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transaction)
	assert.Equal(suite.T(), "2", transaction.ID)
	assert.Equal(suite.T(), suite.testTrx.CustomerID, transaction.CustomerID)
	assert.Equal(suite.T(), suite.testTrx.ActivityType, transaction.ActivityType)
	assert.Equal(suite.T(), suite.testTrx.Amount, transaction.Amount)
	assert.Equal(suite.T(), suite.testTrx.MerchantID, transaction.MerchantID)

	allTransactions, err := suite.repo.FindAllTransaction()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), allTransactions, 2)
}

func (suite *TransactionRepositoryTestSuite) TestFindAllTransaction() {
	transactions, err := suite.repo.FindAllTransaction()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transactions)
	assert.Len(suite.T(), transactions, 1)
	assert.Equal(suite.T(), "1", transactions[0].ID)
	assert.Equal(suite.T(), "1", transactions[0].CustomerID)
	assert.Equal(suite.T(), models.PaymentActivity, transactions[0].ActivityType)
	assert.Equal(suite.T(), 100.0, transactions[0].Amount)
	assert.Equal(suite.T(), "2", transactions[0].MerchantID)

	_, err = suite.repo.CreateTransaction(suite.testTrx)
	assert.NoError(suite.T(), err)

	transactions, err = suite.repo.FindAllTransaction()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transactions)
	assert.Len(suite.T(), transactions, 2)
}

func TestTransactionRepositorySuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}
