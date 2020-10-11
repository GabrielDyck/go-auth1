package internal

import (
	"auth1/pkg/mysql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signUpPath = "/signup"
)



type signupService struct {
	db mysql.SignUp
}


func NewSignUpService(db mysql.SignUp) signupService {
	return signupService{
		db:db,
	}
}

func (s *signupService) signUp(req UserSignReq) error {
    hashedPassword := hashPassword(req.Password)
	return s.db.SignUpAccount(req.Email,hashedPassword,req.AccountType)
}


func SignUp(router *mux.Router, 	client mysql.SignUp) {

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
			WrapBadRequestResponse(writer,err)
		}

	}).Methods("POST")
}

