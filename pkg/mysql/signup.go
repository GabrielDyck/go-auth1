package mysql

import (
	"auth1/api"
	"fmt"
)

type SignUp interface {
	SignUpAccount(email, password string, accountType api.AccountType) error
}


func (c *client) SignUpAccount(email, password string, accountType api.AccountType) error {

	stmt, err := c.db.Prepare(
		"INSERT INTO ACCOUNTS(EMAIL,PASSWORD,ACCOUNT_TYPE,CREATION_DATE)" +
			"VALUES (?,?,?,NOW())")
	if err != nil {
		return err
	}

	result, err := stmt.Exec(email, password, accountType)

	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	fmt.Println(fmt.Sprintf("Created account: %s , type: %s, id: %d", email, accountType, id))
	return nil
}
