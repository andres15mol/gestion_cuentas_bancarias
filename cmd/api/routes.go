package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/", app.ItsAlive)
	// Create an account
	mux.Post("/create-account", app.CreateAccount)
	// Get balance
	mux.Get("/get-balance/{id}", app.GetBalanceById)
	// // Bank deposit
	mux.Post("/bank-deposit", app.BankDepositById)
	// // Bank withdrawal
	mux.Post("/bank-withdrawal", app.BankWithdrawal)
	// // Get transactions
	mux.Get("/get-transactions/{id}", app.getTransactions)

	return mux
}
