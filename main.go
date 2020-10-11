package main

import (
	"auth1/pkg/app"
	"auth1/pkg/config"
	"auth1/pkg/mail"
	"auth1/pkg/routes"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {

	emailSender:= mail.NewSender()
	configuration := config.Read()
	application := app.SetUpApplication(configuration)
	application.Client.Connect()
	var router = mux.NewRouter()
	var authRouter = mux.NewRouter().PathPrefix("/auth").Subrouter()
	customRouter := routes.NewCustomRouter(application.Client,configuration)
	customRouter.AddRoutes(router,configuration.ExpirationDateInMin,emailSender)
	customRouter.AddAuthRoutes(authRouter)


	err := http.ListenAndServe(configuration.Port, nil)
	if err != nil {
		fmt.Println(err)
	}

}
