package main

import (
	"bytes"
	"encoding/json"
	"gestion_cuentas_bancarias/internal/data"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

type testServer struct {
	app application
}

func newTestApplication() *application {
	return &application{
		config:   config{},
		infoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		errorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		models:   data.New(),
	}
}

func newTestServer(t *testing.T, routes http.Handler) *testServer {
	ts := &testServer{
		app: *newTestApplication(),
	}
	return ts
}

func TestItsAlive(t *testing.T) {
	app := newTestApplication()

	// Crear un request HTTP simulado
	req, err := http.NewRequest("GET", "/api/v1/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crear un ResponseRecorder para registrar la respuesta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.ItsAlive)

	// Ejecutar el handler
	handler.ServeHTTP(rr, req)

	// Verificar el código de estado
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verificar el cuerpo de la respuesta
	var response jsonResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Error {
		t.Errorf("Expected Error to be false, got true")
	}

	if response.Message != "The Server is alive" {
		t.Errorf("Expected message 'The Server is alive', got '%s'", response.Message)
	}
}

func TestCreateAccount(t *testing.T) {
	// Preparar el entorno de prueba
	app := newTestApplication()
	
	// Limpiar la base de datos simulada
	data.SimulateAccountDB = make(map[int]data.Account)


	// Crear el payload para la solicitud
	payload := struct {
		OpeningBalance float64 `json:"opening_balance"`
	}{
		OpeningBalance: 1000.0,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	// Crear un request HTTP simulado
	req, err := http.NewRequest("POST", "/api/v1/account", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Crear un ResponseRecorder para registrar la respuesta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreateAccount)

	// Ejecutar el handler
	handler.ServeHTTP(rr, req)

	// Verificar el código de estado
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verificar el cuerpo de la respuesta
	var response jsonResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Error {
		t.Errorf("Expected Error to be false, got true")
	}

	if response.Message != "Created account" {
		t.Errorf("Expected message 'Created account', got '%s'", response.Message)
	}

	// Verificar que los datos de la cuenta existen en la respuesta
	data, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Errorf("Expected Data to be of type map[string]interface{}")
	}

	_, accountExists := data["account_number"]
	if !accountExists {
		t.Errorf("Expected account_number in response data")
	}

	openingBalance, balanceExists := data["Opening balance"]
	if !balanceExists {
		t.Errorf("Expected Opening balance in response data")
	}

	// Verificar que el saldo inicial es correcto
	if openingBalance != float64(1000.0) {
		t.Errorf("Expected Opening balance to be 1000.0, got %v", openingBalance)
	}
}

func TestCreateAccountInvalidJSON(t *testing.T) {
	// Preparar el entorno de prueba
	app := newTestApplication()

	// Crear un JSON inválido
	invalidJSON := []byte(`{"opening_balance": invalid}`)

	// Crear un request HTTP simulado
	req, err := http.NewRequest("POST", "/api/v1/account", bytes.NewBuffer(invalidJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Crear un ResponseRecorder para registrar la respuesta
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.CreateAccount)

	// Ejecutar el handler
	handler.ServeHTTP(rr, req)

	// Verificar el código de estado
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Verificar el cuerpo de la respuesta
	var response jsonResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if !response.Error {
		t.Errorf("Expected Error to be true, got false")
	}

	if response.Message != "invalid json supplied, or json missing entirely" {
		t.Errorf("Expected error message about invalid JSON, got '%s'", response.Message)
	}
}

func TestGetBalanceById(t *testing.T) {
	// Preparar el entorno de prueba
	app := newTestApplication()
	
	// Limpiar la base de datos simulada
	data.SimulateAccountDB = make(map[int]data.Account)

	// Crear una cuenta de prueba
	testAccount := data.Account{
		ID:             20250001,
		OpeningBalance: 500.0,
		LastBalance:    750.0, // Saldo diferente al inicial para probar que se recupera correctamente
		CreatedAt:      time.Now(),
		Transactions:   []int{},
	}
	
	data.SimulateAccountDB[testAccount.ID] = testAccount

	// Crear un router Chi para manejar los parámetros de URL
	r := chi.NewRouter()
	r.Get("/api/v1/account/{id}", app.GetBalanceById)

	// Crear un request HTTP simulado
	req, err := http.NewRequest("GET", "/api/v1/account/20250001", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crear un ResponseRecorder para registrar la respuesta
	rr := httptest.NewRecorder()

	// Ejecutar el handler a través del router
	r.ServeHTTP(rr, req)

	// Verificar el código de estado
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verificar el cuerpo de la respuesta
	var response jsonResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Error {
		t.Errorf("Expected Error to be false, got true")
	}

	if response.Message != "Your balance is" {
		t.Errorf("Expected message 'Your balance is', got '%s'", response.Message)
	}

	// Verificar que los datos de la cuenta existen en la respuesta
	responseData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Errorf("Expected Data to be of type map[string]interface{}")
	}

	accountNumber, accountExists := responseData["account_number"]
	if !accountExists {
		t.Errorf("Expected account_number in response data")
	}

	if float64(accountNumber.(float64)) != float64(testAccount.ID) {
		t.Errorf("Expected account_number to be %d, got %v", testAccount.ID, accountNumber)
	}

	balance, balanceExists := responseData["Balance"]
	if !balanceExists {
		t.Errorf("Expected Balance in response data")
	}

	if balance != float64(750.0) {
		t.Errorf("Expected Balance to be 750.0, got %v", balance)
	}
}

func TestGetBalanceByIdInvalidAccount(t *testing.T) {
	// Preparar el entorno de prueba
	app := newTestApplication()
	
	// Limpiar la base de datos simulada
	data.SimulateAccountDB = make(map[int]data.Account)

	// Crear un router Chi para manejar los parámetros de URL
	r := chi.NewRouter()
	r.Get("/api/v1/account/{id}", app.GetBalanceById)

	// Crear un request HTTP simulado con un ID que no existe
	req, err := http.NewRequest("GET", "/api/v1/account/999999", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Crear un ResponseRecorder para registrar la respuesta
	rr := httptest.NewRecorder()

	// Ejecutar el handler a través del router
	r.ServeHTTP(rr, req)

	// Verificar el código de estado
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	// Verificar el cuerpo de la respuesta
	var response jsonResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if !response.Error {
		t.Errorf("Expected Error to be true, got false")
	}
}
