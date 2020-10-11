package internal

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	logoutPath = "/logout"
)

func Logout(router *mux.Router) {
	router.HandleFunc(logoutPath, func(writer http.ResponseWriter, request *http.Request) {

		user := request.Header.Get("User")
		token := request.Header.Get("Authorization")
		err := validateRequiredHeaders(user, token)

		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
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