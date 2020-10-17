package auth

import "auth1/pkg/mysql"

type AuthService interface {
	IsAuthorized(token string) (bool, error)
	IsProfileAuthorized(id int64, token string) (bool, error)
}

type authService struct {
	db mysql.Auth
}

func NewAuthService(db mysql.Auth) AuthService {
	return &authService{db: db}
}

func (a *authService) IsAuthorized(token string) (bool, error) {
	return a.db.IsAuthenticated(token)
}

func (a *authService) IsProfileAuthorized(id int64, token string) (bool, error) {
	return a.db.IsProfileAuthorized(id, token)
}
