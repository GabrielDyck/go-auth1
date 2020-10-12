package mysql

import (
	"auth1/api"
	"auth1/pkg/mysql/model"
	"fmt"
)

type SignUp interface {
	SignUpBasicAccount(email, password string) error
	SignUpGoogleAccount(email string) error
	AccountAlreadyExists(email string) (bool,error)
	CreateAuthorizationToken(id int64, token string) error
	GetProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*model.Account, error)

}


func (c *client) SignUpBasicAccount(email, password string) error {

	stmt, err := c.db.Prepare(
		"INSERT INTO ACCOUNTS(EMAIL,PASSWORD,ACCOUNT_TYPE)" +
			"VALUES (?,?,\"BASIC\")")
	if err != nil {
		return err
	}
	defer stmt.Close()

	fmt.Println(stmt)
	result, err := stmt.Exec(email, password)

	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	fmt.Println(fmt.Sprintf("Created account: %s , type: BASIC, id: %d", email, id))
	return nil
}

func (c *client) SignUpGoogleAccount(email string) error {

	stmt, err := c.db.Prepare(
		"INSERT INTO ACCOUNTS(EMAIL,ACCOUNT_TYPE)" +
				"VALUES (?,\"GOOGLE\")")

	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(email)

	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	fmt.Println(fmt.Sprintf("Created account: %s , type: GOOGLE, id: %d", email, id))
	return nil
}


func (c *client) AccountAlreadyExists(email string) (bool, error) {
	row, err := c.db.Query("SELECT COUNT(1) FROM ACCOUNTS WHERE EMAIL = ?;", email)

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