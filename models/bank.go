package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	Email            string  `json:"email"`
	Password         string  `json:"password"`
	Address          string  `json:"address"`
	AccountNumber    int     `json:"account_number"`
	AvailableBalance float64 `json:"available_balance"`
}

type Transaction struct {
	gorm.Model
	TransactionType        string  `json:"transaction_type"`
	SenderAccountNumber    int     `json:"sender_account_number"`
	RecipientAccountNumber int     `json:"recipient_account_number"`
	Amount                 float64 `json:"amount"`
	Narration              string  `json:"narration"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AddMoney struct {
	AccountNumber int     `json:"account_number"`
	Amount        float64 `json:"amount"`
	Narration     string  `json:"narration"`
}
