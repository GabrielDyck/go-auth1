package mysql

import (
	"auth1/pkg/mysql/model"
	"fmt"
)

type ForgotPassword interface {
	GetProfileInfoByEmailAndAccountType(email string, accountType model.AccountType) (*model.Account, error)
	DeleteForgotPasswordToken(token int64) error
	CreateForgotPasswordToken(id int64, token string, expirationDateInMIn int) error
	GetForgotPasswordToken(token string)(*model.ForgotPasswordToken,error)
}

func (c *client) CreateForgotPasswordToken(id int64, token string, expirationDateInMIn int) error {

	stmt, err := c.db.Prepare(
		"INSERT INTO FORGOT_PASSWORD_TOKENS(TOKEN,ACCOUNT_ID,EXPIRATION_DATE)" +
			"VALUES (?,?,DATEADD(NOW(),INTERVAL ? MINUTE)")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(token, id, expirationDateInMIn)

	if err != nil {
		return err
	}
	_, _ = result.LastInsertId()
	fmt.Println(fmt.Sprintf("Created forgotPasswordToken: %s, accountID: %d", token, id))
	return nil
}

func (c *client) DeleteForgotPasswordToken(token int64) error {
	_, err := c.db.Exec("DELETE FROM FORGOT_PASSWORD_TOKENS WHERE TOKEN = ?;", token)

	if err != nil {
		return err
	}
	return nil
}


func (c *client) GetForgotPasswordToken(token string)(*model.ForgotPasswordToken,error) {
	row, err := c.db.Query("SELECT TOKEN, ACCOUNT_ID, EXPIRATION_DATE FROM ACCOUNTS WHERE TOKEN = ?;",token)

	if err != nil {
		return nil, err
	}
	var forgotPasswordToken model.ForgotPasswordToken
	if !row.Next() {
		return nil, nil
	}
	err = row.Scan(&forgotPasswordToken.Token, &forgotPasswordToken.AccountID, &forgotPasswordToken.ExpirationDate)

	if err != nil {
		return nil, err
	}

	return &forgotPasswordToken, nil
}
