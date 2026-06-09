package handlers

import (
	"FTBank/middleware"
	"FTBank/models"
	"FTBank/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type BankHandler struct {
	BankService service.BankService
}

func NewBankHandler(service service.BankService) *BankHandler {
	return &BankHandler{BankService: service}
}

func (h BankHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// collect registration details
	var request models.User
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call the service layer
	err = h.BankService.RegisterUser(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "user created successfully")

}

func (h BankHandler) Login(w http.ResponseWriter, r *http.Request) {
	// create a variable to hold the decoded req body
	var loginRequest models.LoginRequest

	// decode req body
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send the details to the service layer
	token, err := h.BankService.Login(loginRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return token as response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

func (h *BankHandler) AddMoney(w http.ResponseWriter, r *http.Request) {
	// decode req body
	var addMoney models.AddMoney
	err := json.NewDecoder(r.Body).Decode(&addMoney)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if user is logged in
	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call service layer
	err = h.BankService.AddMoney(addMoney, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "money added successfully")

}

func (h *BankHandler) TransferMoney(w http.ResponseWriter, r *http.Request) {
	// decode req body
	var transferRequest models.AddMoney
	err := json.NewDecoder(r.Body).Decode(&transferRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// call service layer
	err = h.BankService.TransferMoney(transferRequest, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
