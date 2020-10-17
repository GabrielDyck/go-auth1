package singin

import (
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signInPath = "/signin"
)

func (s *signInService) AddRoutes(router *mux.Router) {

	router.HandleFunc(signInPath,
		func(writer http.ResponseWriter, request *http.Request) {

			resp, status, token := s.signIn(request)
			if token != "" {
				writer.Header().Set("Authorization", token)
			}
			writer.WriteHeader(status)
			writer.Write(resp)

		}).Methods("POST")

}
