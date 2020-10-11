package routes

import (
	"auth1/pkg/mysql"
	"auth1/pkg/mysql/model"
	"crypto/rand"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signInPath = "/signin"
)

type signInService struct {
	db                  mysql.SignIn
	expirationDateInMin int
}

func NewSignInService(db mysql.SignIn, expirationDateInMin int) signInService {
	return signInService{
		db:                  db,
		expirationDateInMin: expirationDateInMin,
	}
}

func (s *signInService) signIn(req UserSignReq) (bool, error) {
	encrypterPassword := hashPassword(req.Password)

	return s.db.IsLoginGranted(req.Email, encrypterPassword)
}

func (s *signInService) getProfileInfo(req UserSignReq) (*model.Account, error) {
	return s.db.GetProfileInfoByEmailAndAccountType(req.Email, req.AccountType)
}

func (s *signInService) generateSessionToken(id int64, expirationDateInMin int) (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenString := string(token)
	err = s.db.CreateAuthorizationToken(id, tokenString, expirationDateInMin)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func signIn(router *mux.Router, db mysql.SignIn, expirationDateInMin int) {

	service := NewSignInService(db, expirationDateInMin)
	router.HandleFunc(signInPath,
		func(writer http.ResponseWriter, request *http.Request) {

			var req UserSignReq
			err := parseRequest(writer, request, &req)
			if err != nil {
				return
			}
			fmt.Println(req)
			_, err = service.signIn(req)

			if err != nil {
				wrapBadRequestResponse(writer, err)
			}

			profileInfo, err := service.getProfileInfo(req)
			if err != nil {
				builtResponse(writer, http.StatusInternalServerError)
				return
			}

			token, err :=service.generateSessionToken(profileInfo.ID,service.expirationDateInMin)

			writer.Header().Set("AUTHORIZATION",token)
			data, httpStatus := builtResponse(profileInfo, http.StatusOK)
			wrapResponse(writer, data, httpStatus)

		}).Methods("POST")
}
