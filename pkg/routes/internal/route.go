package internal

import (
	"auth1/pkg/routes/internal/auth"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	authenticated               = "/authenticated"
)

func Authenticated(router *mux.Router, service auth.AuthService) {

	router.HandleFunc(authenticated, func(w http.ResponseWriter, request *http.Request) {


		token:= request.Header.Get("Authorization")
		isAuthenticated,err:=service.IsAuthorized(token)

		if err !=nil{
			WrapInternalErrorResponse(w,err)
			return
		}

		type AlreadySignIn struct {
			 Authenticated bool `json:"authenticated"`
		}


		data, httpStatus := BuiltResponse(AlreadySignIn{
			Authenticated: isAuthenticated,
		}, http.StatusOK)
		WrapResponse(w,data,httpStatus)	}).Methods("GET")

}