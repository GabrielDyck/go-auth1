package mysql

import (
	"auth1/api"
	"fmt"
	"log"
)

type SignIn interface {
	IsBasicLoginGranted(email, password string) (bool, error)
	GetProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error)
	CreateAuthorizationToken(id int64, token string) error
}



func (c *client) IsBasicLoginGranted(email, password string) (bool, error) {
	row, err := c.db.Query("SELECT COUNT(1) FROM ACCOUNTS WHERE EMAIL = ? AND ACCOUNT_TYPE =\"BASIC\" AND PASSWORD= ?;", email, password)

	if err != nil {
		return false, err
	}
	var count int
	if !row.Next() {
		return false, nil
	}
	err = row.Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 1, nil

}


func (c *client) CreateAuthorizationToken(id int64, token string) error {

	stmt, err := c.db.Prepare(
		"INSERT INTO SESSION_TOKENS(TOKEN,ACCOUNT_ID)" +
			"VALUES (?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token,id)

	if err != nil {
		return err
	}
	log.Println(fmt.Sprintf("Created session token for accountId: %s , token: %s", id, token))
	return nil
}
