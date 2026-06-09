package service

import (
	"FTBank/models"
	"FTBank/repository"
	"FTBank/utils"
	"fmt"
)

type BankService struct {
	BankRepo repository.BankRepository
}

func NewBankService(repo repository.BankRepository) *BankService {
	return &BankService{BankRepo: repo}
}

func (s BankService) RegisterUser(request models.User) error {
	// check if the user already exists
	err := s.BankRepo.CheckUserExists(request.Email)
	if err != nil {
		return fmt.Errorf("user already exist", err)
	}

	// hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	// add the hashed password into the request
	request.Password = hashedPassword

	accountNumber, err := utils.GenerateAccountNumber()
	if err != nil {
		return err
	}

	// add the account number to the request
	request.AccountNumber = accountNumber

	// set the user's account balance to zero
	request.AvailableBalance = 0.0

	// add the user into the db
	err = s.BankRepo.CreateUser(request)
	if err != nil {
		return err
	}

	return nil

}
