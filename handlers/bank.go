package handlers

import (
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
