package mysql

import (
	"auth1/pkg/mysql/model"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Client interface {
	SignUp
	SignIn
	ProfileInfo
	Connect()
}

type SignUp interface {
	SignUpAccount(email, password, accountType string) error
}

type SignIn interface {
	IsLoginGranted(email, password string) (bool, error)
	GetProfileInfoByEmailAndAccountType(email, accountType string) (*model.Account, error)
	CreateAuthorizationToken(id int64, token string, expirationMinutes int) error}

type ProfileInfo interface {
	GetProfileInfoById(id int64) (*model.Account, error)
}

type client struct {
	address  string
	schema   string
	username string

	db *sql.DB
}

func NewClient(address, schema, username string) Client {
	return &client{
		address:  address,
		schema:   schema,
		username: username,
	}
}

func (c *client) Connect() {

	db, err := sql.Open("mysql", c.builtDatasourceName())
	if err != nil {
		panic(fmt.Sprintf("couldn't open mysql connection. %v", err))
	}
	c.db = db
}

func (c *client) IsLoginGranted(email, password string) (bool, error) {
	row, err := c.db.Query("SELECT COUNT(1) FROM ACCOUNTS WHERE ACCOUNT_TYPE !=\"GOOGLE\" AND EMAIL = ? AND PASSWORD= ?;", email, password)

	if err != nil {
		return false, err
	}
	var count int
	err = row.Scan(&count)

	if err != nil {
		return false, err
	}

	return count == 1, nil

}

func (c *client) GetProfileInfoByEmailAndAccountType(email, accountType string) (*model.Account, error) {
	row, err := c.db.Query("SELECT ID, EMAIL, FULLNAME, ADDRESS, PHONE FROM ACCOUNTS WHERE EMAIL = ?  AND ACCOUNT_TYPE !=?;", email, accountType)

	if err != nil {
		return nil, err
	}
	var account model.Account
	err = row.Scan(&account.ID, &account.Email, &account.Fullname, &account.Address, &account.Phone)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *client) GetProfileInfoById(id int64) (*model.Account, error) {
	row, err := c.db.Query("SELECT ID, EMAIL, FULLNAME, ADDRESS, PHONE FROM ACCOUNTS WHERE ID = ?;", id)

	if err != nil {
		return nil, err
	}
	var account model.Account
	err = row.Scan(&account.ID, &account.Email, &account.Fullname, &account.Address, &account.Phone)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (c *client) SignUpAccount(email, password, accountType string) error {

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

func (c *client) CreateAuthorizationToken(id int64, token string, expirationMinutes int) error {

	stmt, err := c.db.Prepare(
		"INSERT INTO SESSION_TOKENS(TOKEN,ACCOUNT_ID,EXPIRATION_DATE)" +
			"VALUES (?,?,DATE_ADD(NOW(), INTERVAL ? MINUTE))")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(id, token, expirationMinutes)

	if err != nil {
		return err
	}
	fmt.Println(fmt.Sprintf("Created session token for accountId: %s , token: %s, expirationMinutes: %d", id, token, expirationMinutes))
	return nil
}

func (c *client) builtDatasourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", c.username, os.Getenv("MYSQL_PASS"), c.address, c.schema)
}
