package singin

import (
	"auth1/api"
	"auth1/pkg/oauth"
	"auth1/pkg/routes/internal"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	signInPath = "/signin"
)

func SignIn(router *mux.Router, service signInService) {

	router.HandleFunc(signInPath,
		func(writer http.ResponseWriter, request *http.Request) {

			var req api.UserSignReq
			err := internal.ParseRequest(writer, request, &req)
			if err != nil {
				return
			}
			log.Println(req)

			var account *api.Account
			switch req.AccountType {

			case api.Basic:
				isGranted, err := service.signIn(req.Email, req.Password)

				if err != nil {
					internal.WrapBadRequestResponse(writer, err)
					return
				}

				if !isGranted {
					internal.WrapBadRequestResponse(writer, errors.New("username or password are wrong"))
					return
				}

				account, err = service.getAccountByEmailAndAccountType(req.Email, api.Basic)
				if err != nil {
					internal.WrapInternalErrorResponse(writer, err)
					return
				}

			case api.Google:

				tokenInfo, err := oauth.VerifyIdToken(req.GoogleToken)
				if err != nil {
					internal.WrapBadRequestResponse(writer, err)
					return
				}

				account, err = service.getAccountByEmailAndAccountType(tokenInfo.Email, api.Google)
				if err != nil {
					internal.WrapInternalErrorResponse(writer, err)
					return
				}

				if account == nil {
					internal.WrapBadRequestResponse(writer, errors.New("user doesn't exists"))
					return

				}

			default:
				internal.WrapBadRequestResponse(writer, errors.New("unknown account type"))
				return

			}
			token, err := service.generateSessionToken(account.ID)
			if err != nil {
				internal.WrapInternalErrorResponse(writer, err)
				return
			}
			writer.Header().Set("Authorization", token)
			data, httpStatus := internal.BuiltResponse(account, http.StatusOK)
			internal.WrapResponse(writer, data, httpStatus)

		}).Methods("POST")

}