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
	application := app.SetUpApplication(configuration)
	application.Client.Connect()
	var router = mux.NewRouter()
	routes.AddRoutes(router,application.Client)


	err := http.ListenAndServe(configuration.Port, router)
	if err != nil {
		fmt.Println(err)
	}

}
