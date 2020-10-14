package internal

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	resetPasswordPath = "/reset-password"
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

		token := request.Header.Get("Forgot-Token")

		var req api.ResetPasswordReq
		err := parseRequest(writer, request, &req)
		if err != nil {
			return
		}
		log.Println(token)

		forgotPasswordToken, err := resetPasswordService.db.GetForgotPasswordTokenByToken(token)
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}

		if forgotPasswordToken == nil {
			WrapBadRequestResponse(writer, errors.New("token doesn't exist in our database"))
			return
		}

		//TODO FIX THIS
		splitedTime:=strings.Split(time.Now().String(),".")
		now,_:= time.Parse("2006-01-02 15:04:05",splitedTime[0])
		if forgotPasswordToken.ExpirationDate.Before(now)  {
			WrapBadRequestResponse(writer, errors.New("token has expired"))
			return
		}
		account, err := resetPasswordService.db.GetAccountById(forgotPasswordToken.AccountID)
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}

		err = resetPasswordService.db.ChangePassword(account.ID, hashPassword(req.Password))
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}
		err = resetPasswordService.db.DeleteForgotPasswordToken(token)
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}

		WrapOkEmptyResponse(writer)
	}).Methods("POST")
}
