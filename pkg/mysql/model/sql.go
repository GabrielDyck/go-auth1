package model


type AccountType string

const(
	Basic  AccountType = "BASIC"
	Google AccountType = "GOOGLE"
)
type Account struct {
	ID          int64
	Email       string
	Fullname    string
	Address     string
	AccountType AccountType
	Phone       string
}


type ForgotPasswordToken struct {
	AccountID int64
	Token string
	ExpirationDate string
}