package routes

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

const (
	signUp = "/signup"
)

type SignUpReq struct {
	Email string `json:"email"`
	Password string`json:"password"`
}

func addSignUp(router *mux.Router) {
	router.HandleFunc(signUp, func(writer http.ResponseWriter, request *http.Request) {

		var req SignUpReq
		err := parseRequest(writer, request,&req)
		if err!=nil {
			return
		}
		fmt.Println(req)

	}).Methods("POST")
}

func parseRequest(writer http.ResponseWriter, request *http.Request, bodyStruct interface {}) error {
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		wrapBadRequest(writer, err)
		return err
	}

	err = json.Unmarshal(body, bodyStruct)

	if err != nil {
		wrapBadRequest(writer, err)
	}
	return err
}

func wrapBadRequest(writer http.ResponseWriter, err error) {
	data, httpStatus := writeResponse(builtErrorBodyMsg(err), http.StatusBadRequest)
	writer.WriteHeader(httpStatus)
	_,err = writer.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}
