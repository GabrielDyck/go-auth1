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

func (l *logoutService) Logout(router *mux.Router) {
	router.HandleFunc(logoutPath, func(writer http.ResponseWriter, request *http.Request) {

		token := request.Header.Get("Authorization")
		err := validateRequiredHeaders(token)

		if err != nil {
			internal.WrapBadRequestResponse(writer, err)
			return
		}

		ok,err:= l.db.Logout(token)

		if err != nil {
			internal.WrapInternalErrorResponse(writer, err)
			return
		}

		if !ok{
			internal.WrapBadRequestResponse(writer, errors.New("token isn't available"))
			return
		}

		internal.WrapOkEmptyResponse(writer)



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