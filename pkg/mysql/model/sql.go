package model

import (
	"auth1/api"
	"time"
)


type Account struct {
	ID          int64
	Email       string
	Fullname    string
	Address     string
	AccountType api.AccountType
	Phone       string
}


type ForgotPasswordToken struct {
	AccountID int64
	Token string
	ExpirationDate time.Time
}