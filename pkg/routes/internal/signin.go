package internal

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/oauth"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signInPath = "/signin"
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
	encrypterPassword := hashPassword(password)

	return s.db.IsLoginGranted(email, encrypterPassword)
}

func (s *signInService) getAccountByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error) {
	return s.db.GetProfileInfoByEmailAndAccountType(email, accountType)
}

func (s *signInService) generateSessionToken(id int64) (string, error) {
	token := make([]byte, 255)
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

func SignIn(router *mux.Router, service signInService) {

	router.HandleFunc(signInPath,
		func(writer http.ResponseWriter, request *http.Request) {

			var req api.UserSignReq
			err := parseRequest(writer, request, &req)
			if err != nil {
				return
			}
			fmt.Println(req)

			var account *api.Account
			switch req.AccountType {

			case api.Basic:
				isGranted, err := service.signIn(req.Email, req.Password)

				if err != nil {
					WrapBadRequestResponse(writer, err)
					return
				}

				if !isGranted {
					WrapBadRequestResponse(writer, errors.New("username or password are wrong"))
					return
				}

				account, err = service.getAccountByEmailAndAccountType(req.Email, api.Basic)
				if err != nil {
					WrapInternalErrorResponse(writer, err)
					return
				}

			case api.Google:

				tokenInfo, err := oauth.VerifyIdToken(req.GoogleToken)
				if err != nil {
					WrapBadRequestResponse(writer, err)
					return
				}

				account, err = service.getAccountByEmailAndAccountType(tokenInfo.Email, api.Google)
				if err != nil {
					WrapInternalErrorResponse(writer, err)
					return
				}

				if account == nil {
					WrapBadRequestResponse(writer, errors.New("user doesn't exists"))
					return

				}

			default:
				WrapBadRequestResponse(writer, errors.New("unknown account type"))
				return

			}
			token, err := service.generateSessionToken(account.ID)
			writer.Header().Set("AUTHORIZATION", token)
			data, httpStatus := builtResponse(account, http.StatusOK)
			wrapResponse(writer, data, httpStatus)

		}).Methods("POST")

}
