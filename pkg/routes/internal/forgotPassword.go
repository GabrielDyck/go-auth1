package internal

import (
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
	"auth1/pkg/mysql/model"
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

type ForgotPasswordReq struct {
	Email string `json:"email"`
}

func NewForgotPasswordService( db mysql.ForgotPassword)forgotPasswordService{
	return forgotPasswordService{
		db: db,
	}
}

func ForgotPassword(router *mux.Router, db mysql.ForgotPassword, expirationDateInMin int,emailSender mail.Sender) {
	service := NewForgotPasswordService(db)
	router.HandleFunc(forgotPasswordPath, func(writer http.ResponseWriter, request *http.Request) {

		var req ForgotPasswordReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)

		account, err :=service.db.GetProfileInfoByEmailAndAccountType(req.Email,model.Basic)
		if err != nil {
			WrapBadRequestResponse(writer, err)
		}

		if account== nil{
			WrapBadRequestResponse(writer, errors.New("email doesn't exist in our database"))
		}

		err =service.db.CreateForgotPasswordToken(account.ID,expirationDateInMin)

		if err != nil {
			WrapInternalErrorResponse(writer, err)
		}

		emailSender.SendEmail(req.Email)

		writer.WriteHeader(http.StatusOK)

	}).Methods("POST")
}