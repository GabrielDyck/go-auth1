package auth

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	authenticated = "/authenticated"
)

func (a *authService) AddRoutes(router *mux.Router) {

	router.HandleFunc(authenticated, func(writer http.ResponseWriter, request *http.Request) {

		resp, status := a.Authenticated(request)
		writer.WriteHeader(status)
		writer.Write(resp)

	}).Methods("GET")

}
