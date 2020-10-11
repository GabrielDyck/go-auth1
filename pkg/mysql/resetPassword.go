package mysql

import (
	"auth1/pkg/mysql/model"
	"time"
)

type ResetPassword interface {
	DeleteForgotPasswordToken(token string) error
	GetForgotPasswordTokenByToken(token string) (*model.ForgotPasswordToken, error)
	ChangePassword(accountId int64,newPassword string) error
	Account
}


func (c *client) DeleteForgotPasswordToken(token string) error {
	_, err := c.db.Exec("DELETE FROM FORGOT_PASSWORD_TOKENS WHERE TOKEN = ?;", token)

	if err != nil {
		return err
	}
	return nil
}

func (c *client) GetForgotPasswordTokenByToken(token string) (*model.ForgotPasswordToken, error) {
	row, err := c.db.Query("SELECT TOKEN, ACCOUNT_ID, EXPIRATION_DATE FROM ACCOUNTS WHERE TOKEN = ?;", token)

	if err != nil {
		return nil, err
	}
	var forgotPasswordToken model.ForgotPasswordToken
	if !row.Next() {
		return nil, nil
	}

	var sqlDateTime string

	err = row.Scan(&forgotPasswordToken.Token, &forgotPasswordToken.AccountID, &sqlDateTime)

	if err != nil {
		return nil, err
	}

	forgotPasswordToken.ExpirationDate, err = time.Parse(c.datetimeLayout, sqlDateTime)

	if err != nil {
		return nil, err
	}

	return &forgotPasswordToken, nil
}



func (c *client) ChangePassword(accountId int64,newPassword string)error {
	_, err := c.db.Exec("UPDATE FROM ACCOUNTS SET PASSWORD= ? WHERE ID = ?;",newPassword, accountId)

	if err != nil {
		return err
	}
	return nil
}