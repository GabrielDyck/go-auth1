package mysql

import (
	"auth1/pkg/mysql/model"
	"context"
	"database/sql"
	"time"
)

type ResetPassword interface {
	DeleteForgotPasswordToken(tx *sql.Tx, ctx context.Context,token string) error
	GetForgotPasswordTokenByToken(token string) (*model.ForgotPasswordToken, error)
	ChangePassword(tx *sql.Tx, ctx context.Context,accountId int64,newPassword string) error
	CreateTrx(context context.Context)(*sql.Tx, error)
	Account
}


func (c *client) DeleteForgotPasswordToken(tx *sql.Tx, ctx context.Context,token string) error {
	_, err := tx.ExecContext(ctx,"DELETE FROM FORGOT_PASSWORD_TOKENS WHERE TOKEN = ?;", token)

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (c *client) GetForgotPasswordTokenByToken(token string) (*model.ForgotPasswordToken, error) {
	row, err := c.db.Query("SELECT TOKEN, ACCOUNT_ID, EXPIRATION_DATE FROM FORGOT_PASSWORD_TOKENS WHERE TOKEN = ?;", token)

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



func (c *client) ChangePassword(tx *sql.Tx, ctx context.Context,accountId int64,newPassword string)error {
	_, err := tx.ExecContext(ctx,"UPDATE ACCOUNTS SET PASSWORD= ? WHERE ID = ?;",newPassword, accountId)

	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}