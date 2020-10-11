package internal

import (
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
	"auth1/pkg/mysql/model"
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

func ForgotPassword(router *mux.Router,  db mysql.ForgotPassword, emailSender mail.Sender) {
	service := NewForgotPasswordService(db)
	router.HandleFunc(forgotPasswordPath, func(writer http.ResponseWriter, request *http.Request) {

		var req ForgotPasswordReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)

		service.db.GetProfileInfoByEmailAndAccountType(req.Email,model.Basic)

	}).Methods("POST")
}
