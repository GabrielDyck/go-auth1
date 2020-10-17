package forgot

import (
	"auth1/api"
	"auth1/pkg/mail"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/internal"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	forgotPasswordPath = "/forgot-password"
)

type forgotPasswordService struct {
	db                  mysql.ForgotPassword
	expirationDateInMin int
	emailSender         mail.Sender
}

func NewForgotPasswordService(db mysql.ForgotPassword, expirationDateInMin int, emailSender mail.Sender) forgotPasswordService {
	return forgotPasswordService{
		db:                  db,
		expirationDateInMin: expirationDateInMin,
		emailSender:         emailSender,
	}
}

func (f *forgotPasswordService) AddRoutes(router *mux.Router){
	router.HandleFunc(forgotPasswordPath, func(writer http.ResponseWriter, request *http.Request) {
		resp,status:= f.forgotPassword(request)
		writer.WriteHeader(status)
		writer.Write(resp)
	}).Methods("POST")
}


func (f *forgotPasswordService) forgotPassword(request *http.Request) ([]byte, int) {

	var req api.ForgotPasswordReq
	err := internal.ParseRequest(request, &req)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}
	log.Println(req)

	account, err := f.db.GetProfileInfoByEmailAndAccountType(req.Email, api.Basic)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}

	if account == nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("email doesn't exist in our database")), http.StatusBadRequest)

	}

	token := f.tokenGenerator()
	err = f.db.CreateForgotPasswordToken(account.ID, f.expirationDateInMin, token)

	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)

	}

	f.emailSender.SendEmail(req.Email, token)

	return []byte("{}"), http.StatusOK

}

func (f *forgotPasswordService) tokenGenerator() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%X", b)

}
