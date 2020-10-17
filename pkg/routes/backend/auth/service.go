package auth

import (
	"auth1/pkg/mysql"
	"auth1/pkg/routes/internal"
	"github.com/gorilla/mux"
	"net/http"
)

type AuthService interface {
	IsAuthorized(token string) (bool, error)
	IsProfileEditorAuthorized(id int64, token string) (bool, error)
	AddRoutes(router *mux.Router)
}

type authService struct {
	db mysql.Auth
}

func NewAuthService(db mysql.Auth) AuthService {
	return &authService{db: db}
}

func (a *authService) Authenticated(request *http.Request) ([]byte, int) {

	token := request.Header.Get("Authorization")
	isAuthenticated, err := a.IsAuthorized(token)

	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}

	type AlreadySignIn struct {
		Authenticated bool `json:"authenticated"`
	}

	return internal.BuiltResponse(AlreadySignIn{
		Authenticated: isAuthenticated,
	}, http.StatusOK)

}

func (a *authService) IsAuthorized(token string) (bool, error) {
	return a.db.IsAuthenticated(token)
}

func (a *authService) IsProfileEditorAuthorized(id int64, token string) (bool, error) {
	return a.db.IsProfileAuthorized(id, token)
}
