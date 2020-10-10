package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	resetPassword = "/resetPassword"
)

type ResetPasswordReq struct {
	Password string `json:"password"`
}

func addResetPassword(router *mux.Router) {
	router.HandleFunc(resetPassword, func(writer http.ResponseWriter, request *http.Request) {

		token :=request.Header.Get("FORGOT_TOKEN")

		var req ForgotPasswordReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(token)

	}).Methods("POST")
}
