package model

import (
	"time"
)


type ForgotPasswordToken struct {
	AccountID int64
	Token string
	ExpirationDate time.Time
}