package repository

import (
	"FTBank/models"

	"gorm.io/gorm"
)

type BankRepository interface {
	CheckUserExists(email string) (models.User, error)
	CreateUser(request models.User) error
	GetUserByID(userID uint) (models.User, error)
	UpdateUser(user models.User) error
	GetUserByAccountNumber(accountNumber int) (models.User, error)
	MakeTransfer(sender, recipient models.User, transferRequest models.AddMoney) error
}

type PostgresBankRepo struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) BankRepository {
	return &PostgresBankRepo{db: db}
}

func (r PostgresBankRepo) CheckUserExists(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r PostgresBankRepo) CreateUser(request models.User) error {
	err := r.db.Create(&request).Error
	if err != nil {
		return err
	}

	return nil
}

func (r PostgresBankRepo) GetUserByID(userID uint) (models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *PostgresBankRepo) UpdateUser(user models.User) error {
	err := r.db.Save(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresBankRepo) GetUserByAccountNumber(accountNumber int) (models.User, error) {
	var user models.User
	err := r.db.Where("account_number = ?", accountNumber).First(&user).Error
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *PostgresBankRepo) MakeTransfer(sender, recipient models.User, transferRequest models.AddMoney) error {
	txn := r.db.Begin()

	// deduct from sender's balance
	sender.AvailableBalance -= transferRequest.Amount

	// add the amount to the recipient
	recipient.AvailableBalance += transferRequest.Amount

	// save the transaction for sender
	err := txn.Save(sender).Error
	if err != nil {
		txn.Rollback()
		return err
	}

	// save the transaction for recipient
	err = txn.Save(recipient).Error
	if err != nil {
		txn.Rollback()
		return err
	}

	// record the transaction
	transaction := &models.Transaction{
		TransactionType:        "transfer",
		SenderAccountNumber:    sender.AccountNumber,
		RecipientAccountNumber: recipient.AccountNumber,
		Amount:                 transferRequest.Amount,
		Narration:              transferRequest.Narration,
	}

	// 
	err = txn.Create(transaction).Error
	if err != nil {
		txn.Rollback()
		return err
	}

	txn.Commit()

	return nil

}
