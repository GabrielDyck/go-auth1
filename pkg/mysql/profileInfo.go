package mysql

import "auth1/pkg/mysql/model"

type Account interface {
	GetAccountById(id int64) (*model.Account, error)
}


func (c *client) GetAccountById(id int64) (*model.Account, error) {
	row, err := c.db.Query("SELECT ID, EMAIL, FULLNAME, ADDRESS, PHONE FROM ACCOUNTS WHERE ID = ?;", id)

	if err != nil {
		return nil, err
	}
	var account model.Account
	if !row.Next() {
		return nil, nil
	}
	err = row.Scan(&account.ID, &account.Email, &account.Fullname, &account.Address, &account.Phone)

	if err != nil {
		return nil, err
	}

	return &account, nil
}
