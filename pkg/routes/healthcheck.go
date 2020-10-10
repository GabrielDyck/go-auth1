package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	healthCheck = "/health-check"
)

func addHealthCheck(router *mux.Router) {
	router.HandleFunc(healthCheck, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_,err:= writer.Write([]byte("OK"))
		if err != nil {
			fmt.Println(err)
		}
	}).Methods("GET")
}
