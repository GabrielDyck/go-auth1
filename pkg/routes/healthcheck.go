package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	healthCheckPath = "/health-check"
)

func healthCheck(router *mux.Router) {
	router.HandleFunc(healthCheckPath, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_,err:= writer.Write([]byte("OK"))
		if err != nil {
			fmt.Println(err)
		}
	}).Methods("GET")
}
