package routes

import (
	"auth1/pkg/mysql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signUpPath = "/signup"
)

type Service interface {
	SignUp(req UserSignReq) bool
}


type Response struct {
	Email  string          `json:"email"`
	Result OperationResult `json:"result"`
}

type service struct {
	db mysql.SignUp
}


func NewSignUpService(db mysql.SignUp) service {
	return service{
		db:db,
	}
}

func (s *service) signUp(req UserSignReq) error {
	return s.db.SignUpAccount(req.Email,req.Password,req.AccountType)
}


func addSignUp(router *mux.Router, 	client mysql.SignUp) {

	service:= NewSignUpService(client)
	router.HandleFunc(signUpPath, func(writer http.ResponseWriter, request *http.Request) {

		var req UserSignReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)
		err = service.signUp(req)

		if err !=nil{
			wrapBadRequestResponse(writer,err)
		}

	}).Methods("POST")
}

