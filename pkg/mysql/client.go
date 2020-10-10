package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Client interface {
	SignUp
	Connect()
}

type SignUp interface {
	SignUpAccount(username, password, accountType string) error
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

func (c *client) SignUpAccount(username , password, accountType string) error{

	stmt,err:= c.db.Prepare(
		"INSERT INTO ACCOUNTS(USERNAME,PASSWORD,ACCOUNT_TYPE,CREATION_DATE)" +
		"VALUES (?,?,?,NOW())")
	if err != nil{
		return err
	}

	result, err := stmt.Exec(username,password,accountType)

	if err != nil {
		return err
	}
	id,_ := result.LastInsertId()
	fmt.Println(fmt.Sprintf("Created account: %s , type: %s, id: %d",username,accountType,id))
	return nil
}

func (c *client) builtDatasourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", c.username, os.Getenv("MYSQL_PASS"), c.address, c.schema)
}
