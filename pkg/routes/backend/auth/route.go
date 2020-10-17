package auth

import (
	"auth1/pkg/routes/internal"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	authenticated               = "/authenticated"
)

func Authenticated(router *mux.Router, service AuthService) {

	router.HandleFunc(authenticated, func(w http.ResponseWriter, request *http.Request) {


		token:= request.Header.Get("Authorization")
		isAuthenticated,err:=service.IsAuthorized(token)

		if err !=nil{
			internal.WrapInternalErrorResponse(w,err)
			return
		}

		type AlreadySignIn struct {
			 Authenticated bool `json:"authenticated"`
		}


		data, httpStatus := internal.BuiltResponse(AlreadySignIn{
			Authenticated: isAuthenticated,
		}, http.StatusOK)
		internal.WrapResponse(w,data,httpStatus)

	}).Methods("GET")

}