package data

import (
    "testing"
    "github.com/stretchr/testify/assert"
)
// TestInsertAccount tests the Insert method of the Account model with initial balance. 
func TestInsertAccount(t *testing.T) {
    models := New()

    account := Account{
        OpeningBalance: 100.0,
    }

    accountID, err := models.Account.Insert(account)
    assert.NoError(t, err)
    assert.NotZero(t, accountID)

    insertedAccount, err := models.Account.GetByID(accountID)
    assert.NoError(t, err)
    assert.Equal(t, accountID, insertedAccount.ID)
    assert.Equal(t, account.OpeningBalance, insertedAccount.LastBalance)
}

// ********************************************

// TestInsertTransaction tests the InsertTransaction method of the Transaction model.
func TestInsertTransaction(t *testing.T) {
    models := New()

    account := Account{
        OpeningBalance: 100.0,
    }

    accountID, err := models.Account.Insert(account)
    assert.NoError(t, err)

    transactionID, err := models.Transaction.InsertTransaction("Deposit", accountID, 50.0)
    assert.NoError(t, err)
    assert.NotZero(t, transactionID)

    insertedTransaction := SimulateTransactionDB[transactionID]
    assert.Equal(t, transactionID, insertedTransaction.ID)
    assert.Equal(t, "Deposit", insertedTransaction.Type)
    assert.Equal(t, 50.0, insertedTransaction.Amount)
    assert.Equal(t, 150.0, insertedTransaction.Balance)
}


// TestInsertTransactionInvalidAccount tests the InsertTransaction method of the Transaction model with an invalid account.
func TestInsertTransactionInvalidType(t *testing.T) {
    models := New()

    account := Account{
        OpeningBalance: 100.0,
    }

    accountID, err := models.Account.Insert(account)
    assert.NoError(t, err)

    _, err = models.Transaction.InsertTransaction("InvalidType", accountID, 50.0)
    assert.Error(t, err)
    assert.Equal(t, "Invalid transaction type", err.Error())
}

// ********************************************

// TestInsertTransaction_InsufficientFunds tests the InsertTransaction method of the Transaction model with insufficient funds.
func TestInsertTransaction_InsufficientFunds(t *testing.T) {
    models := New()

    account := Account{
        OpeningBalance: 100.0,
    }

    accountID, err := models.Account.Insert(account)
    assert.NoError(t, err)

    _, err = models.Transaction.InsertTransaction("Withdrawal", accountID, 150.0)
    assert.Error(t, err)
    assert.Equal(t, "Insufficient funds", err.Error())
}

// TestGetTransactionsByAccountID tests the GetTransactionsByAccountID method of the Transaction model.
func TestGetTransactionsByAccountID(t *testing.T) {
    models := New()

    account := Account{
        OpeningBalance: 100.0,
    }

    accountID, err := models.Account.Insert(account)
    assert.NoError(t, err)

    _, err = models.Transaction.InsertTransaction("Deposit", accountID, 50.0)
    assert.NoError(t, err)

    transactions, err := models.Transaction.GetTransactionsByAccountID(accountID)
    assert.NoError(t, err)
    assert.Len(t, transactions, 1)
    assert.Equal(t, "Deposit", transactions[0].Type)
    assert.Equal(t, 50.0, transactions[0].Amount)
    assert.Equal(t, 150.0, transactions[0].Balance)
}
