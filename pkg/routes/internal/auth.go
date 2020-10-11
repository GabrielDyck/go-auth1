package internal

import "auth1/pkg/mysql"

type AuthService interface {
	isAuthorized(token string) (bool,error)
}

type authService struct {
	db mysql.Auth
}

func NewAuthService(db mysql.Auth) AuthService {
	return &authService{db: db}
}


func (a *authService) isAuthorized(token string) (bool, error) {
	return a.db.IsAuthenticated(token)
}


