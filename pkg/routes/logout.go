package routes

import (
	"errors"
	"fmt"
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
		err := validateRequiredHeaders(user, token)

		if err != nil {
			wrapBadRequestResponse(writer, err)
		}

	}).Methods("POST")
}

func validateRequiredHeaders(headers ...string) error {

	for _, header := range headers {
		if header == "" {
			return errors.New(fmt.Sprintf("%s is not present", header))
		}
	}
	return nil
}