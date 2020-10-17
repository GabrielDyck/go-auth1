package main

import (
	"auth1/pkg/app"
	"auth1/pkg/config"
	"auth1/pkg/mail"
	"auth1/pkg/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

func main() {

	emailSender:= mail.NewSender()
	configuration := config.Read()
	application := app.SetUpApplication(configuration)
	application.Client.Connect()
	var router = mux.NewRouter().PathPrefix("/backend").Subrouter()
	var authRouter = mux.NewRouter().PathPrefix("/auth").Subrouter()
	customRouter := routes.NewCustomRouter(application.Client,configuration)
	customRouter.AddBackendRoutes(router,configuration.ExpirationDateInMin,emailSender)
	customRouter.AddAuthRoutes(authRouter)
	customRouter.AddFrontendRoutes()


	err := http.ListenAndServe(os.Getenv("PORT"), nil)
	if err != nil {
		log.Println(err)
	}

}
