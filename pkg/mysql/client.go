package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

type Client interface {
	Connect()
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

func (c *client) builtDatasourceName() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", c.username, os.Getenv("MYSQL_PASS"), c.address, c.schema)
}
