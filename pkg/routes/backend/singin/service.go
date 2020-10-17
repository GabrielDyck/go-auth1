package singin

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/internal"
	"crypto/rand"
	"fmt"
)

type signInService struct {
	db mysql.SignIn
}

func NewSignInService(db mysql.SignIn) signInService {
	return signInService{
		db: db,
	}
}

func (s *signInService) signIn(email, password string) (bool, error) {
	encrypterPassword := internal.HashPassword(password)

	return s.db.IsBasicLoginGranted(email, encrypterPassword)
}

func (s *signInService) getAccountByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error) {
	return s.db.GetProfileInfoByEmailAndAccountType(email, accountType)
}

func (s *signInService) generateSessionToken(id int64) (string, error) {
	token := make([]byte, 128)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenString := fmt.Sprintf("%X", token)
	err = s.db.CreateAuthorizationToken(id, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}
