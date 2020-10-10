package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signUp = "/signup"
)


func addSignUp(router *mux.Router) {
	router.HandleFunc(signUp, func(writer http.ResponseWriter, request *http.Request) {

		var req UserSignReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)

	}).Methods("POST")
}
