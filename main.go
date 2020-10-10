package main

import (
	"auth1/pkg/config"
	"auth1/pkg/routes"
	"fmt"
	"net/http"
)

func main(){


	configuration := config.Read()

	routes.AddRoutes()
	err := http.ListenAndServe(configuration.Port,http.DefaultServeMux)
	if err !=nil {
		fmt.Println(err)
	}

}