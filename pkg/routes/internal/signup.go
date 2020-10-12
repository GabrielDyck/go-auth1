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
	signUpPath               = "/signup"
	accountAlreadyExistsPath = "/accountAlreadyExistsPath"
)

type signupService struct {
	db mysql.SignUp
}

func NewSignUpService(db mysql.SignUp) signupService {
	return signupService{
		db: db,
	}
}

func (s *signupService) signUpBasicAccount(req api.UserSignReq) error {
	fmt.Println("signUpBasicAccount")

	hashedPassword := hashPassword(req.Password)
	return s.db.SignUpBasicAccount(req.Email, hashedPassword)
}

// TODO : extract with signin
func (s *signupService) generateSessionToken(id int64) (string, error) {
	fmt.Println("generateSessionToken")

	token := make([]byte, 128)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenString := fmt.Sprintf("%X",token)
	return s.createAuthToken(id, tokenString)
}

func (s *signupService) createAuthToken(id int64, tokenString string) (string, error) {
	err := s.db.CreateAuthorizationToken(id, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *signupService) signUpGoogleAccount(email string) error {
	fmt.Println("SignUpGoogleAccount")
	return s.db.SignUpGoogleAccount(email)
}
func (s *signupService) accountAlreadyExists(email string) (bool, error) {
	fmt.Println("AccountAlreadyExists")

	return s.db.AccountAlreadyExists(email)
}

func (s *signupService) getProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error) {
	fmt.Println("GetProfileInfoByEmailAndAccountType")
	return s.db.GetProfileInfoByEmailAndAccountType(email, accountType)
}

func SignUp(router *mux.Router, client mysql.SignUp) {

	service := NewSignUpService(client)
	router.HandleFunc(signUpPath, func(writer http.ResponseWriter, request *http.Request) {

		var req api.UserSignReq
		err := parseRequest(writer, request, &req)
		if err != nil {
			return
		}
		fmt.Println(req)

		var account *api.Account

		switch req.AccountType {

		case api.Basic:
			already, err := service.accountAlreadyExists(req.Email)

			if err != nil {
				WrapInternalErrorResponse(writer, err)
				return
			}

			if already {
				WrapBadRequestResponse(writer, errors.New("user already registered"))
				return
			}

			err = service.signUpBasicAccount(req)
			if err != nil {
				WrapInternalErrorResponse(writer, err)
				return
			}

			account, err = service.getProfileInfoByEmailAndAccountType(req.Email, api.Basic)
			if err != nil {
				WrapInternalErrorResponse(writer, err)
				return
			}

		case api.Google:
			tokenInfo,err := oauth.VerifyIdToken(req.GoogleToken)
			if err != nil {
				WrapBadRequestResponse(writer, err)
				return
			}

			account, err = service.getProfileInfoByEmailAndAccountType(tokenInfo.Email, api.Google)
			if err != nil {
				WrapInternalErrorResponse(writer, err)
				return
			}

			if account == nil {
				err=service.signUpGoogleAccount(tokenInfo.Email)
				if err != nil {
					WrapInternalErrorResponse(writer, err)
					return
				}
			}
		default:
			WrapBadRequestResponse(writer, errors.New("unknown account type"))
			return
		}
		WrapOkEmptyResponse(writer)
	}).Methods("POST")

	router.HandleFunc(accountAlreadyExistsPath, func(writer http.ResponseWriter, request *http.Request) {

		var req api.UserSignReq
		err := parseRequest(writer, request, &req)
		if err != nil {
			return
		}
		fmt.Println(req)

		alreadyExists, err := service.accountAlreadyExists(req.Email)
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		data, httpStatus := builtResponse(alreadyExists, http.StatusOK)
		wrapResponse(writer, data, httpStatus)
	}).Methods("POST")
}
