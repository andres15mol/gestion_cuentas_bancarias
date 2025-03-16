package main

import (
	"gestion_cuentas_bancarias/internal/data"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

var staticPath = "./static/"

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type envelope map[string]interface {
}

func (app *application) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		OpeningBalance float64 `json:"opening_balance"`
	}
	var payload jsonResponse

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	var Account data.Account

	Account.OpeningBalance = requestPayload.OpeningBalance

	Account.ID, err = app.models.Account.Insert(Account)

	//send back a response
	payload = jsonResponse{
		Error:   false,
		Message: "Created account",
		Data:    envelope{"account_number": Account.ID, "Opening balance": Account.OpeningBalance},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) GetBalanceById(w http.ResponseWriter, r *http.Request) {
	var payload jsonResponse

	accountID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	account, err := app.models.Account.GetByID(accountID)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = err.Error()
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	//send back a response
	payload = jsonResponse{
		Error:   false,
		Message: "Your balance is",
		Data: envelope{"account_number": account.ID,
			"Balance": account.LastBalance},
	}

	// out, err := json.MarshalIndent(payload, "", "\t")
	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}

}

func (app *application) BankDepositById(w http.ResponseWriter, r *http.Request) {

	app.makeTransaction(w, r, "Deposit")

}

func (app *application) BankWithdrawal(w http.ResponseWriter, r *http.Request) {

	app.makeTransaction(w, r, "Withdrawal")

}

func (app *application) makeTransaction(w http.ResponseWriter, r *http.Request, transactionType string) {
	var requestPayload struct {
		AccountID int     `json:"account_id"`
		Amount    float64 `json:"amount"`
	}
	var payload jsonResponse

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	//Insert transaction on the DB
	transactionID, err := app.models.Transaction.InsertTransaction(transactionType, requestPayload.AccountID, requestPayload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = err.Error()
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	// Message
	var message string
	if transactionType == "Deposit" {
		message = "Deposit successful"
	} else {
		message = "Withdrawal successful"
	}

	//Send back a response
	payload = jsonResponse{
		Error:   false,
		Message: message,
		Data:    envelope{"transaction_id": transactionID, "account_id": requestPayload.AccountID, "amount": requestPayload.Amount},
	}

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) getTransactions(w http.ResponseWriter, r *http.Request) {
	var payload jsonResponse

	accountID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	transactions, err := app.models.Transaction.GetTransactionsByAccountID(accountID)
	if err != nil {
		app.errorLog.Println(err)
		payload.Error = true
		payload.Message = err.Error()
		_ = app.writeJSON(w, http.StatusBadRequest, payload)
		return
	}

	//send back a response
	payload = jsonResponse{
		Error:   false,
		Message: "Your transactions are",
		Data:    envelope{"Final Balance": transactions[len(transactions)-1].Balance , "Transactions": transactions},
	}
	

	err = app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ItsAlive(w http.ResponseWriter, r *http.Request) {

	//send back a response
	payload := jsonResponse{
		Error:   false,
		Message: "The Server is alive",
	}

	// out, err := json.MarshalIndent(payload, "", "\t")
	err := app.writeJSON(w, http.StatusOK, payload)
	if err != nil {
		app.errorLog.Println(err)
	}

}
