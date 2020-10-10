package routes

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	logout = "/logout"
)

func addLogout(router *mux.Router) {
	router.HandleFunc(logout, func(writer http.ResponseWriter, request *http.Request) {


		user := request.Header.Get("User")
		token := request.Header.Get("Authorization")
		err:= validateRequiredHeaders(user, token)

		if err != nil {
			wrapBadRequest(writer,err)
		}

	}).Methods("POST")
}

func validateRequiredHeaders(user, token string) error {
	if len(user) == 0 {
		return errors.New("user is not present")
	} else if len(token) == 0 {
		return errors.New("token is not present")
	}

	return nil
}
