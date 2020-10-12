package internal

import (
	"auth1/api"
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	forgotPasswordPath = "/forgotPassword"
)

type forgotPasswordService struct {
	db mysql.ForgotPassword
}



func NewForgotPasswordService( db mysql.ForgotPassword)forgotPasswordService{
	return forgotPasswordService{
		db: db,
	}
}

func ForgotPassword(router *mux.Router, db mysql.ForgotPassword, expirationDateInMin int,emailSender mail.Sender) {
	service := NewForgotPasswordService(db)
	router.HandleFunc(forgotPasswordPath, func(writer http.ResponseWriter, request *http.Request) {

		var req api.ForgotPasswordReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)

		account, err :=service.db.GetProfileInfoByEmailAndAccountType(req.Email,api.Basic)
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}

		if account== nil{
			WrapBadRequestResponse(writer, errors.New("email doesn't exist in our database"))
			return
		}

		token := service.tokenGenerator()
		err =service.db.CreateForgotPasswordToken(account.ID,expirationDateInMin,token)

		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}

		emailSender.SendEmail(req.Email,token)

		writer.WriteHeader(http.StatusOK)

	}).Methods("POST")
}

func (f * forgotPasswordService) tokenGenerator() string {
	b := make([]byte, 32)
	_,_=rand.Read(b)
	return fmt.Sprintf("%x", b)
}