package api

import "time"

type AccountType string

const(
	Basic  AccountType = "BASIC"
	Google AccountType = "GOOGLE"
)



type ErrorMSG struct {
	Reason string `json:"reason"`
}

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

type ResetPasswordReq struct {
	Password string `json:"password"`
}


type UserSignReq struct {
	Email       string      `json:"email"`
	Password    string      `json:"password"`
	GoogleToken string      `json:"token"`
	AccountType AccountType `json:"account_type"`
}

type UserSessionTokenResponse struct {
	SessionToken string `json:"session_token"`
}


type Account struct {
	ID          int64 `json:"id"`
	Email       string`json:"email"`
	Fullname    *string`json:"fullname"`
	Address     *string`json:"address"`
	AccountType AccountType`json:"account_type"`
	Phone       *string`json:"phone"`
}


type ForgotPasswordToken struct {
	AccountID int64
	Token string
	ExpirationDate time.Time
}