package routes

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)


type OperationResult string

const(
	Error OperationResult = "ERROR"
	Success OperationResult = "SUCCESS"
)
type UserSignReq struct {
	Email string `json:"email"`
	Password string`json:"password"`
	AccountType string `json:"account_type"`
}

type ErrorMSG struct {
	Reason string `json:"reason"`
}

func builtResponse(response interface{}, statusCode int)([]byte ,int) {

	data , err :=json.Marshal(response)

	if err !=nil {
		data, err=json.Marshal(response)

		if err != nil {
			fmt.Println(fmt.Sprintf("error occurred trying to marshal error response: %v",err))
		}
		return data, http.StatusInternalServerError
	}

	return data,statusCode
}


func builtErrorBodyMsg(err error) ErrorMSG{
	return ErrorMSG{Reason: err.Error()}
}



func parseRequest(writer http.ResponseWriter, request *http.Request, bodyStruct interface {}) error {
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		wrapBadRequestResponse(writer, err)
		return err
	}

	err = json.Unmarshal(body, bodyStruct)

	if err != nil {
		wrapBadRequestResponse(writer, err)
	}
	return err
}

func wrapResponse(writer http.ResponseWriter,data []byte, httpStatus int) {
	writer.WriteHeader(httpStatus)
	_,err := writer.Write(data)
	if err != nil {
		fmt.Println(err)
	}
}


func wrapInternalErrorResponse(writer http.ResponseWriter, err error) {
	data, httpStatus := builtResponse(builtErrorBodyMsg(err), http.StatusInternalServerError)
	wrapResponse(writer,data,httpStatus)

}

func wrapBadRequestResponse(writer http.ResponseWriter, err error) {
	data, httpStatus := builtResponse(builtErrorBodyMsg(err), http.StatusBadRequest)
	wrapResponse(writer,data,httpStatus)

}

func hashPassword(pass string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pass))
	encrypterPassword := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return encrypterPassword
}


