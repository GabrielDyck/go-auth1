package logout

import (
	"auth1/pkg/mysql"
	"auth1/pkg/routes/internal"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	logoutPath = "/logout"
)

type logoutService struct {
	db mysql.Logout
}

func NewLogoutService(db mysql.Logout) logoutService {
	return logoutService{
		db: db,
	}
}

func (l *logoutService) logout(request *http.Request) ([]byte, int) {

	token := request.Header.Get("Authorization")
	err := validateRequiredHeaders(token)

	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}

	ok, err := l.db.Logout(token)

	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}

	if !ok {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("token isn't available")), http.StatusBadRequest)
	}
	return  []byte("{}"), http.StatusOK

}

func (l *logoutService) AddRoutes(router *mux.Router) {
	router.HandleFunc(logoutPath, func(writer http.ResponseWriter, request *http.Request) {
		resp, status := l.logout(request)
		writer.WriteHeader(status)
		writer.Write(resp)

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
