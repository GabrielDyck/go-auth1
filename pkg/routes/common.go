package routes

import (
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

func writeResponse(response interface{}, statusCode int)([]byte ,int) {

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

