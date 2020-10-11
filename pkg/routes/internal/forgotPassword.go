package internal

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	forgotPasswordPath = "/forgotPassword"
)

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

func ForgotPassword(router *mux.Router) {
	router.HandleFunc(forgotPasswordPath, func(writer http.ResponseWriter, request *http.Request) {

		var req ForgotPasswordReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)

	}).Methods("POST")
}
