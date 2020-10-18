package mysql

import (
	"auth1/api"
)

type Account interface {
	GetAccountById(id int64) (*api.Account, error)
	AccountAlreadyExists(email string, accountType api.AccountType) (bool, error)
	EditProfileInfo(accountId int64,email,address,fullname,phone string) error
}


func (c *client) GetAccountById(id int64) (*api.Account, error) {
	row, err := c.db.Query("SELECT ID, EMAIL, FULLNAME, ADDRESS, PHONE, ACCOUNT_TYPE FROM ACCOUNTS WHERE ID = ?;", id)

	if err != nil {
		return nil, err
	}
	defer row.Close()

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


func (c *client) EditProfileInfo(accountId int64,email,address,fullname,phone string) error {
	_, err := c.db.Exec("UPDATE ACCOUNTS SET EMAIL= ?, ADDRESS= ?, FULLNAME=?, PHONE=? WHERE ID = ?;",email, address, fullname, phone, accountId)

	if err != nil {
		return err
	}

	return nil
}
