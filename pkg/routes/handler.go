package routes

import (
	"fmt"
	"net/http"
)

const (
	healthCheck = "/health-check"
)
func AddRoutes(){
	addHealthCheck()
}




func addHealthCheck(){
	http.HandleFunc(healthCheck, func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		_,err:= writer.Write([]byte("OK"))
		if err != nil {
			fmt.Println(err)
		}
	})
}