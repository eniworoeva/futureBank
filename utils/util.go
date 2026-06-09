package utils

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const (
	min = 11111111
	max = 99999999
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func ComparePasswords(plainPass, HashedPass string) error {
	err := bcrypt.CompareHashAndPassword([]byte(HashedPass), []byte(plainPass))
	if err != nil {
		return err
	}

	return nil
}

func GenerateAccountNumber() (int, error) {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min, nil
}