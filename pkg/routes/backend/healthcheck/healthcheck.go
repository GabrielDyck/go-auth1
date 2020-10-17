package healthcheck

import (
	"github.com/gorilla/mux"
	"log"
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
			log.Println(err)
		}
	}).Methods("GET")
}
