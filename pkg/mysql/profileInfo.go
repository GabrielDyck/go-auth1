package mysql

import (
	"auth1/api"
)

type Account interface {
	GetAccountById(id int64) (*api.Account, error)
}


func (c *client) GetAccountById(id int64) (*api.Account, error) {
	row, err := c.db.Query("SELECT ID, EMAIL, FULLNAME, ADDRESS, PHONE, ACCOUNT_TYPE FROM ACCOUNTS WHERE ID = ?;", id)

	if err != nil {
		return nil, err
	}
	var account api.Account
	if !row.Next() {
		return nil, nil
	}
	err = row.Scan(&account.ID, &account.Email, &account.Fullname, &account.Address, &account.Phone, &account.AccountType)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
