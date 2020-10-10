package main

import (
	"auth1/pkg/app"
	"auth1/pkg/config"
	"auth1/pkg/routes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	configuration := config.Read()
	app := app.SetUpApplication(configuration)
	app.Client.Connect()
	var router = mux.NewRouter()
	routes.AddRoutes(router)


	err := http.ListenAndServe(configuration.Port, router)
	if err != nil {
		fmt.Println(err)
	}

}
