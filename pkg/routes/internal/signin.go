package internal

import (
	"auth1/api"
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
}

func NewSignInService(db mysql.SignIn) signInService {
	return signInService{
		db:                  db,
	}
}

func (s *signInService) signIn(req api.UserSignReq) (bool, error) {
	encrypterPassword := hashPassword(req.Password)

	return s.db.IsLoginGranted(req.Email, encrypterPassword)
}

func (s *signInService) getProfileInfo(req api.UserSignReq) (*model.Account, error) {
	return s.db.GetProfileInfoByEmailAndAccountType(req.Email, req.AccountType)
}

func (s *signInService) generateSessionToken(id int64) (string, error) {
	token := make([]byte, 16)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenString := string(token)
	err = s.db.CreateAuthorizationToken(id, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func SignIn(router *mux.Router, db mysql.SignIn) {

	service := NewSignInService(db)
	router.HandleFunc(signInPath,
		func(writer http.ResponseWriter, request *http.Request) {

			var req api.UserSignReq
			err := parseRequest(writer, request, &req)
			if err != nil {
				return
			}
			fmt.Println(req)
			_, err = service.signIn(req)

			if err != nil {
				WrapBadRequestResponse(writer, err)
			}

			profileInfo, err := service.getProfileInfo(req)
			if err != nil {
				builtResponse(writer, http.StatusInternalServerError)
				return
			}

			token, err :=service.generateSessionToken(profileInfo.ID)

			writer.Header().Set("AUTHORIZATION",token)
			data, httpStatus := builtResponse(profileInfo, http.StatusOK)
			wrapResponse(writer, data, httpStatus)

		}).Methods("POST")
}
