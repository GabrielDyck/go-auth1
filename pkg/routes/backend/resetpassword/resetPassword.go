package resetpassword

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/internal"
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

func (r *resetPasswordService) resetPassword(request *http.Request) ([]byte, int) {
	token := request.Header.Get("Forgot-Token")

	var req api.ResetPasswordReq
	err := internal.ParseRequest(request, &req)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)

	}
	log.Println(token)

	forgotPasswordToken, err := r.db.GetForgotPasswordTokenByToken(token)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}

	if forgotPasswordToken == nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("token doesn't exist in our database")), http.StatusBadRequest)
	}

	splitedTime := strings.Split(time.Now().String(), ".")
	now, _ := time.Parse("2006-01-02 15:04:05", splitedTime[0])
	if forgotPasswordToken.ExpirationDate.Before(now) {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("token has expired")), http.StatusBadRequest)
	}
	account, err := r.db.GetAccountById(forgotPasswordToken.AccountID)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("token has expired")), http.StatusInternalServerError)
	}

	err = r.db.ChangePassword(account.ID, internal.HashPassword(req.Password))
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}
	err = r.db.DeleteForgotPasswordToken(token)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}

	return []byte{}, http.StatusOK
}

func (r *resetPasswordService) AddRoutes(router *mux.Router) {

	router.HandleFunc(resetPasswordPath, func(writer http.ResponseWriter, request *http.Request) {
		resp, status := r.resetPassword(request)
		writer.WriteHeader(status)
		writer.Write(resp)
	}).Methods("POST")
}
