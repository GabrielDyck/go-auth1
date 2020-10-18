package mysql

import (
	"auth1/api"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"sync"
	"time"
)

type Client interface {
	Auth
	SignUp
	SignIn
	Account
	ResetPassword
	ForgotPassword
	Logout
	Connect()
	//TODO Shutdown()
}

var once sync.Once
var instance Client

type client struct {
	address  string
	schema   string
	username string
	datetimeLayout string
	maxConnection int
	maxIdleConnection int

	db *sql.DB
}

func NewClient(address, schema, username string, maxConnection ,maxIdleConnection int) Client {
	once.Do(func() {
		instance = &client{
			address:  address,
			schema:   schema,
			username: username,
			maxConnection: maxConnection,
			maxIdleConnection: maxIdleConnection,
			datetimeLayout: "2006-01-02 15:04:05",

		}
	})

	return instance

}



func (c *client) GetProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error) {
	row, err := c.db.Query("SELECT ID, EMAIL, FULLNAME, ADDRESS, ACCOUNT_TYPE,PHONE FROM ACCOUNTS WHERE EMAIL = ?  AND ACCOUNT_TYPE =?;", email, accountType)

	if err != nil {
		return nil, err
	}
	defer row.Close()

	var account api.Account
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
	db.SetMaxIdleConns(c.maxIdleConnection)
	db.SetMaxOpenConns(c.maxConnection)

	db.SetConnMaxIdleTime(time.Duration(60)* time.Second)
	c.db = db
}

func (c *client) CreateTrx(context context.Context)(*sql.Tx, error)  {
	return c.db.BeginTx(context, nil)
}



func (c *client) builtDatasourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", c.username, os.Getenv("MYSQL_PASS"), c.address, c.schema)
}
