package internal

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

const (
	resetPasswordPath = "/resetPassword"
)


type resetPasswordService struct {
	db mysql.ResetPassword
}

func NewResetPasswordService(db mysql.ResetPassword) resetPasswordService {
	return resetPasswordService{
		db: db,
	}
}

func ResetPassword(router *mux.Router, db mysql.ResetPassword) {

	resetPasswordService := NewResetPasswordService(db)

	router.HandleFunc(resetPasswordPath, func(writer http.ResponseWriter, request *http.Request) {

		token := request.Header.Get("FORGOT_TOKEN")

		var req api.ResetPasswordReq
		err := parseRequest(writer, request, &req)
		if err != nil {
			return
		}
		fmt.Println(token)

		forgotPasswordToken, err := resetPasswordService.db.GetForgotPasswordTokenByToken(token)
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}

		if forgotPasswordToken == nil {
			WrapBadRequestResponse(writer, errors.New("token doesn't exist in our database"))
			return
		}

		if forgotPasswordToken.ExpirationDate.Unix() < time.Now().Unix() {
			WrapBadRequestResponse(writer, errors.New("token has expired"))
			return
		}
		account, err := resetPasswordService.db.GetAccountById(forgotPasswordToken.AccountID)
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}

		err = resetPasswordService.db.ChangePassword(account.ID, req.Password)
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}
		err = resetPasswordService.db.DeleteForgotPasswordToken(token)
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}

	}).Methods("POST")
}
