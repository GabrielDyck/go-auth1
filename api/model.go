package api

type AccountType string

const(
	Basic  AccountType = "BASIC"
	Google AccountType = "GOOGLE"
)

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

type ResetPasswordReq struct {
	Password string `json:"password"`
}


type UserSignReq struct {
	Email string `json:"email"`
	Password string`json:"password"`
	AccountType AccountType `json:"account_type"`
}
