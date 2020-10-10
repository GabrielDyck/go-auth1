package routes

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	profile = "/profile/{id}"
)

type ProfileWriteReq struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

func addProfileRoutes(router *mux.Router) {
	router.HandleFunc(profile, func(writer http.ResponseWriter, request *http.Request) {
		id := mux.Vars(request)["id"]

		fmt.Println(id)

	}).Methods("GET")

	router.HandleFunc(profile, func(writer http.ResponseWriter, request *http.Request) {
		id := mux.Vars(request)["id"]

		var req ProfileWriteReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(id)
		fmt.Println(req)

		err=validateRequiredFields(req)
		if err !=nil{
			wrapBadRequest(writer,err)
			return
		}

	}).Methods("POST")

}



func validateRequiredFields(req ProfileWriteReq) error{

	if req.Email == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}
