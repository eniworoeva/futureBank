package repository

import (
	"FTBank/models"

	"gorm.io/gorm"
)

type BankRepository interface {
	CheckUserExists(email string) error
	CreateUser(request models.User) error
}

type PostgresBankRepo struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) BankRepository {
	return &PostgresBankRepo{db: db}
}

func (r PostgresBankRepo) CheckUserExists(email string) error {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	// if error is zero, it means we found someone(that is the error)
	if err == nil {
		return err
	}

	return nil
}

func (r PostgresBankRepo) CreateUser(request models.User) error {
	err := r.db.Create(&request).Error
	if err != nil {
		return err
	}

	return nil
}
