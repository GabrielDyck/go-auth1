package internal

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/mysql/model"
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
	err = s.db.CreateAuthorizationToken(id, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *signupService) signBasicUpGoogleAccount(req api.UserSignReq) error {
	fmt.Println("SignUpGoogleAccount")
	return s.db.SignUpGoogleAccount(req.Email)
}
func (s *signupService) accountAlreadyExists(email string) (bool, error) {
	fmt.Println("AccountAlreadyExists")

	return s.db.AccountAlreadyExists(email)
}

func (s *signupService) getProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*model.Account, error) {
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

		var account *model.Account

		switch req.AccountType {

		case api.Basic:
			already, err := service.accountAlreadyExists(req.Email)

			if err != nil {
				WrapInternalErrorResponse(writer, err)
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
			err = service.signBasicUpGoogleAccount(req)
			account, err = service.getProfileInfoByEmailAndAccountType(req.Email, api.Basic)
			if err != nil {
				WrapInternalErrorResponse(writer, err)
				return
			}

			if account ==nil {
				err=service.signBasicUpGoogleAccount(req)
				if err != nil {
					WrapInternalErrorResponse(writer, err)
					return
				}
			}


			token,err:= service.generateSessionToken(account.ID)
			if err != nil {
				WrapInternalErrorResponse(writer, err)
				return
			}

			writer.Header().Set("AUTHORIZATION",token)

		default:
			WrapBadRequestResponse(writer, errors.New("unknown account type"))
			return

		}

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
