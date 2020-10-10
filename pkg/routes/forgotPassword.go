package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	forgotPassword = "/forgotPassword"
)

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

func addForgotPassword(router *mux.Router) {
	router.HandleFunc(forgotPassword, func(writer http.ResponseWriter, request *http.Request) {

		var req ForgotPasswordReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)

	}).Methods("POST")
}
