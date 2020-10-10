package app

import (
	"auth1/pkg/config"
	"auth1/pkg/mysql"
)

type Application struct {
	Client mysql.Client
}

func SetUpApplication(configuration config.Configuration) *Application {

	return &Application{
		Client: mysql.NewClient(configuration.Mysql.Address, configuration.Mysql.Schema, configuration.Mysql.Username),
	}

}
