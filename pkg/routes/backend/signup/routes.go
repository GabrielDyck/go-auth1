package signup

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signUpPath               = "/signup"
)


func (s *signupService) AddRoutes(router *mux.Router) {

	router.HandleFunc(signUpPath, func(writer http.ResponseWriter, request *http.Request) {
		data,status,token:= s.signUp(request)
		if token != ""{
			writer.Header().Set("Authorization", token)
		}
		writer.WriteHeader(status)
		writer.Write(data)


	}).Methods("POST")

}
