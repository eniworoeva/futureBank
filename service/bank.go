package service

import (
	"FTBank/middleware"
	"FTBank/models"
	"FTBank/repository"
	"FTBank/utils"
	"fmt"
	"strconv"
)

type BankService struct {
	BankRepo repository.BankRepository
}

func NewBankService(repo repository.BankRepository) *BankService {
	return &BankService{BankRepo: repo}
}

func (s BankService) RegisterUser(request models.User) error {
	// check if the user already exists
	_, err := s.BankRepo.CheckUserExists(request.Email)
	if err == nil {
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

func (s BankService) Login(loginRequest models.LoginRequest) (string, error) {
	// check if the user exists
	user, err := s.BankRepo.CheckUserExists(loginRequest.Email)
	if err != nil {
		return "", err
	}

	// compare passwords to make sure they match
	err = utils.ComparePasswords(loginRequest.Password, user.Password)
	if err != nil {
		return "", err
	}

	//convert userID fromn uint -> int -> string
	userID := strconv.Itoa(int(user.ID))

	// generate the token
	token, err := middleware.GenerateJWT(userID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *BankService) AddMoney(request models.AddMoney, userID uint) error {
	// search for user
	user, err := s.BankRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// validate the amount
	if request.Amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	// add amount to the user balance
	user.AvailableBalance += request.Amount

	// update user
	err = s.BankRepo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *BankService) TransferMoney(transferRequest models.AddMoney, userID uint) error {
	// validate amount
	if transferRequest.Amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	// check if the recipient exists using their acct number
	recipient, err := s.BankRepo.GetUserByAccountNumber(transferRequest.AccountNumber)
	if err != nil {
		return err
	}

	// search for the sender
	sender, err := s.BankRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// make sure sender have enough in their available balance
	if sender.AvailableBalance < transferRequest.Amount {
		return fmt.Errorf("insufficient funds")
	}

	// make transfer
	err = s.BankRepo.MakeTransfer(sender, recipient, transferRequest)
	if err != nil {
		return err
	}

	return nil
}
