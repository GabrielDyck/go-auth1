package singin

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/oauth"
	"auth1/pkg/routes/internal"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type signInService struct {
	db mysql.SignIn
}

func NewSignInService(db mysql.SignIn) signInService {
	return signInService{
		db: db,
	}
}



func (s *signInService) signIn(request *http.Request) ([]byte, int, string) {

	var req api.UserSignReq
	err := internal.ParseRequest(request, &req)
	if err != nil {
		data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
		return data, status, ""
	}
	log.Println(req)

	var account *api.Account
	switch req.AccountType {

	case api.Basic:
		isGranted, err := s.login(req.Email, req.Password)

		if err != nil {
			data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
			return data, status, ""

		}

		if !isGranted {
			data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("username or password are wrong")), http.StatusBadRequest)
			return data, status, ""
		}

		account, err = s.getAccountByEmailAndAccountType(req.Email, api.Basic)
		if err != nil {
			data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			return data, status, ""
		}

	case api.Google:

		tokenInfo, err := oauth.VerifyIdToken(req.GoogleToken)
		if err != nil {
			data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
			return data, status, ""

		}

		account, err = s.getAccountByEmailAndAccountType(tokenInfo.Email, api.Google)
		if err != nil {
			data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			return data, status, ""

		}

		if account == nil {
			data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("user doesn't exists")), http.StatusBadRequest)
			return data, status, ""

		}

	default:
		data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("unknown account type")), http.StatusBadRequest)
		return data, status, ""

	}
	token, err := s.generateSessionToken(account.ID)
	if err != nil {
		data, status := internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
		return data, status, ""

	}
	data, status := internal.BuiltResponse(account, http.StatusOK)
	return data, status, token
}



func (s *signInService) login(email, password string) (bool, error) {
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
