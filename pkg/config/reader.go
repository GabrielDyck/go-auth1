package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type mysql struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Schema   string `json:"schema"`
}

type Configuration struct {
	Port  string `json:"port"`
	ExpirationDateInMin int `json:"expiration_date_in_min"`
	Mysql mysql  `json:"mysql"`
}

func Read() Configuration {
	var configuration Configuration

	data, err := ioutil.ReadFile("./pkg/config/resources/application.json")

	if err != nil {
		panic(fmt.Sprintf("couldn't read configuration: %v", err))
	}

	err = json.Unmarshal(data, &configuration)
	if err != nil {
		panic(fmt.Sprintf("couldn't unmarshall configuration: %v", err))
	}

	return configuration
}
