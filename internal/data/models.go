package data

import (
	"fmt"
	"sync"
	"time"
)

var (
	mu                    sync.Mutex
	nextID                = 1 // Comienza en 1 o cualquier otro valor
	nextIDTransaction     = 1 // Comienza en 1 o cualquier otro valor
	SimulateAccountDB     = make(map[int]Account)
	SimulateTransactionDB = make(map[int]Transaction)
)

func New() Models {

	return Models{
		Account:     Account{},
		Transaction: Transaction{},
	}
}

type Models struct {
	Account     Account
	Transaction Transaction
}

// Account model
type Account struct {
	ID             int       `json:"id"`
	OpeningBalance float64   `json:"opening_balance"`
	LastBalance    float64   `json:"balance"`
	CreatedAt      time.Time `json:"created_at"`
	Transactions   []int     `json:"transactions"`
}

func generateAccountID() int {
	currentYear := time.Now().Year() // Get Actual Year
	mu.Lock()
	accountID := currentYear*1000000 + nextID // Generate an Id with format Year-000001
	nextID++                                  // Incriese the nextID
	mu.Unlock()                               // Unlock
	return accountID
}

// Insert method to insert a new account
func (a *Account) Insert(account Account) (int, error) {
	account.ID = generateAccountID() // Assign the generated ID to the account
	account.LastBalance = account.OpeningBalance
	account.CreatedAt = time.Now()

	SimulateAccountDB[account.ID] = account

	return account.ID, nil

}

func (a *Account) GetByID(id int) (*Account, error) {

	account, ok := SimulateAccountDB[id]
	if !ok {
		return nil, fmt.Errorf("Account with ID %d not found", id)
	}

	return &account, nil
}

// Transaction model
type Transaction struct {
	ID        int       `json:"id"`
	Amount    float64   `json:"amount"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	Balance   float64   `json:"balance"`
}

func generateTransactionID() int {
	mu.Lock()
	accountID := nextIDTransaction // Generate an Id with format Year-000001
	nextIDTransaction++            // Incriese the nextID
	mu.Unlock()                    // Unlock
	return accountID
}

func (a *Transaction) InsertTransaction(typeTransaction string, accountID int, amount float64) (int, error) {
	// Validate if the account exist
	account, ok := SimulateAccountDB[accountID]
	if !ok {
		return 0, fmt.Errorf("account with ID %d not found", accountID)
	}

	// Create a new transaction
	transaction := Transaction{}

	// Update the account balance

	if typeTransaction == "Deposit" {
		account.LastBalance += amount

	} else if typeTransaction == "Withdrawal" {

		if account.LastBalance < amount {
			return 0, fmt.Errorf("Insufficient funds")
		} else {
			account.LastBalance -= amount
		}

	} else {
		return 0, fmt.Errorf("Invalid transaction type")

	}

	// Assign the values to the transaction
	transaction.ID = generateTransactionID()
	transaction.Type = typeTransaction
	transaction.Amount = amount
	transaction.CreatedAt = time.Now()
	transaction.Balance = account.LastBalance

	// Insert the transaction in the DB
	SimulateTransactionDB[transaction.ID] = transaction

	account.Transactions = append(account.Transactions, transaction.ID)

	// Update the account in the DB
	SimulateAccountDB[accountID] = account

	return transaction.ID, nil

}

func (a *Transaction) GetTransactionsByAccountID(accountID int) ([]Transaction, error) {

	account, ok := SimulateAccountDB[accountID]
	if !ok {
		return nil, fmt.Errorf("Account with ID %d not found", accountID)
	}
	var transactions []Transaction

	if account.Transactions == nil {
		return transactions, fmt.Errorf("Transactions with Account ID %d not found", accountID)

	} else {
		for i := 0; i < len(account.Transactions); i++ {

			transactions = append(transactions, SimulateTransactionDB[account.Transactions[i]])
		}

		return transactions, nil
	}

}
