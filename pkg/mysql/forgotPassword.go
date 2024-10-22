package mysql

import (
	"auth1/api"
	"fmt"
	"log"
)

type ForgotPassword interface {
	GetProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error)
	CreateForgotPasswordToken(id int64, expirationDateInMin int, token string) error
}

func (c *client) CreateForgotPasswordToken(id int64, expirationDateInMin int,token string) error {

	stmt, err := c.db.Prepare(
		"INSERT INTO FORGOT_PASSWORD_TOKENS(TOKEN,ACCOUNT_ID,EXPIRATION_DATE)" +
			"VALUES (?,?,DATE_ADD(NOW(),INTERVAL ? MINUTE))")
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(token, id, expirationDateInMin)

	if err != nil {
		return err
	}
	_, _ = result.LastInsertId()
	log.Println(fmt.Sprintf("Created forgotPasswordToken, accountID: %d", id))
	return nil
}