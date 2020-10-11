package routes

import (
	"auth1/pkg/mysql"
	"auth1/pkg/mysql/model"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	signInPath = "/signin"
)

type signInService struct {
	db mysql.SignIn
}

func NewSignInService(db mysql.SignIn) signInService {
	return signInService{
		db: db,
	}
}

func (s *signInService) signIn(req UserSignReq) (bool, error) {
	return s.db.IsLoginGranted(req.Email, req.Password)
}

func (s *signInService) getProfileInfo(req UserSignReq) (*model.Account, error) {
	return s.db.GetProfileInfoByEmailAndAccountType(req.Email,req.AccountType)
}


func SignIn(router *mux.Router, db mysql.SignIn) {

	service := NewSignInService(db)
	router.HandleFunc(signInPath,
		func(writer http.ResponseWriter, request *http.Request) {

			var req UserSignReq
			err := parseRequest(writer, request, &req)
			if err != nil {
				return
			}
			fmt.Println(req)
			_, err = service.signIn(req)

			if err != nil {
				wrapBadRequestResponse(writer, err)
			}

			profileInfo,err := service.getProfileInfo(req)
			if err != nil {
				builtResponse(writer, http.StatusInternalServerError)
			}

			data, httpStatus :=builtResponse(profileInfo, http.StatusOK)
			wrapResponse(writer,data,httpStatus)
		}).Methods("POST")
}
