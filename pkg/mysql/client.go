package mysql

import (
	"auth1/pkg/mysql/model"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
)

type Client interface {
	SignUp
	SignIn
	ProfileInfo
	Auth
	ForgotPassword
	Connect()
}

var once sync.Once
var instance Client

type client struct {
	address  string
	schema   string
	username string

	db *sql.DB
}

func NewClient(address, schema, username string) Client {
	once.Do(func() {
		instance = &client{
			address:  address,
			schema:   schema,
			username: username,
		}
	})

	return instance

}



func (c *client) GetProfileInfoByEmailAndAccountType(email string, accountType model.AccountType) (*model.Account, error) {
	row, err := c.db.Query("SELECT ID, EMAIL, FULLNAME, ADDRESS, ACCOUNT_TYPE,PHONE FROM ACCOUNTS WHERE EMAIL = ?  AND ACCOUNT_TYPE !=?;", email, accountType)

	if err != nil {
		return nil, err
	}
	var account model.Account
	if !row.Next() {
		return nil, nil
	}
	err = row.Scan(&account.ID, &account.Email, &account.Fullname, &account.Address,&account.AccountType, &account.Phone)

	if err != nil {
		return nil, err
	}

	return &account, nil
}


func (c *client) Connect() {

	db, err := sql.Open("mysql", c.builtDatasourceName())
	if err != nil {
		panic(fmt.Sprintf("couldn't open mysql connection. %v", err))
	}
	c.db = db
}



func (c *client) builtDatasourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", c.username, os.Getenv("MYSQL_PASS"), c.address, c.schema)
}
