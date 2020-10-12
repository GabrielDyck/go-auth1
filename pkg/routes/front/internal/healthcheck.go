package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	healthCheckPath = "/health-check"
)

func HealthCheck(router *mux.Router) {
	router.HandleFunc(healthCheckPath, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_,err:= writer.Write([]byte("OK"))
		if err != nil {
			fmt.Println(err)
		}
	}).Methods("GET")
}
