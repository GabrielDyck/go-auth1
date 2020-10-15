package singup

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/oauth"
	"auth1/pkg/routes/internal"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	signUpPath               = "/signup"
	accountAlreadyExistsPath = "/accountAlreadyExists"
)

func SignUp(router *mux.Router, client mysql.SignUp) {

	service := NewSignUpService(client)
	router.HandleFunc(signUpPath, func(writer http.ResponseWriter, request *http.Request) {

		var req api.UserSignReq
		err := internal.ParseRequest(writer, request, &req)
		if err != nil {
			return
		}
		log.Println(req)

		var account *api.Account

		switch req.AccountType {

		case api.Basic:
			account, err = service.getProfileInfoByEmailAndAccountType(req.Email, api.Basic)

			if err != nil {
				internal.WrapInternalErrorResponse(writer, err)
				return
			}

			if account!=nil {
				internal.WrapBadRequestResponse(writer, errors.New("user already registered"))
				return
			}

			err = service.signUpBasicAccount(req)
			if err != nil {
				internal.WrapInternalErrorResponse(writer, err)
				return
			}

			account, err = service.getProfileInfoByEmailAndAccountType(req.Email, api.Basic)
			if err != nil {
				internal.WrapInternalErrorResponse(writer, err)
				return
			}

		case api.Google:
			tokenInfo,err := oauth.VerifyIdToken(req.GoogleToken)
			if err != nil {
				internal.WrapBadRequestResponse(writer, err)
				return
			}

			account, err = service.getProfileInfoByEmailAndAccountType(tokenInfo.Email, api.Google)
			if err != nil {
				internal.WrapInternalErrorResponse(writer, err)
				return
			}

			if account == nil {
				err=service.signUpGoogleAccount(tokenInfo.Email)
				if err != nil {
					internal.WrapInternalErrorResponse(writer, err)
					return
				}
			}

			account, err = service.getProfileInfoByEmailAndAccountType(tokenInfo.Email, api.Google)
			if err != nil {
				internal.WrapInternalErrorResponse(writer, err)
				return
			}
		default:
			internal.WrapBadRequestResponse(writer, errors.New("unknown account type"))
			return
		}
		data, httpStatus := internal.BuiltResponse(account, http.StatusOK)
		internal.WrapResponse(writer, data, httpStatus)
	}).Methods("POST")

	router.HandleFunc(accountAlreadyExistsPath, func(writer http.ResponseWriter, request *http.Request) {

		var req api.UserSignReq
		err := internal.ParseRequest(writer, request, &req)
		if err != nil {
			return
		}
		log.Println(req)

		alreadyExists, err := service.accountAlreadyExists(req.Email)
		if err != nil {
			internal.WrapBadRequestResponse(writer, err)
			return
		}
		data, httpStatus := internal.BuiltResponse(alreadyExists, http.StatusOK)
		internal.WrapResponse(writer, data, httpStatus)
	}).Methods("POST")
}
