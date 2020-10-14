package internal

import (
	"auth1/api"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)




func builtResponse(response interface{}, statusCode int)([]byte ,int) {

	data , err :=json.Marshal(response)

	if err !=nil {
		data, err=json.Marshal(response)

		if err != nil {
			log.Println(fmt.Sprintf("error occurred trying to marshal error response: %v",err))
		}
		return data, http.StatusInternalServerError
	}

	return data,statusCode
}


func builtErrorBodyMsg(err error) api.ErrorMSG {
	return api.ErrorMSG{Reason: err.Error()}
}



func parseRequest(writer http.ResponseWriter, request *http.Request, bodyStruct interface {}) error {
	body, err := ioutil.ReadAll(request.Body)

	if err != nil {
		WrapBadRequestResponse(writer, err)
		return err
	}

	err = json.Unmarshal(body, bodyStruct)

	if err != nil {
		WrapBadRequestResponse(writer, err)
	}
	return err
}

func wrapResponse(writer http.ResponseWriter,data []byte, httpStatus int) {
	writer.WriteHeader(httpStatus)
	_,err := writer.Write(data)
	if err != nil {
		log.Println(err)
	}
}


func WrapInternalErrorResponse(writer http.ResponseWriter, err error) {
	data, httpStatus := builtResponse(builtErrorBodyMsg(err), http.StatusInternalServerError)
	wrapResponse(writer,data,httpStatus)

}

func WrapBadRequestResponse(writer http.ResponseWriter, err error) {
	data, httpStatus := builtResponse(builtErrorBodyMsg(err), http.StatusBadRequest)
	wrapResponse(writer,data,httpStatus)

}

func WrapNotAllowedRequestResponse(writer http.ResponseWriter, err error) {
	data, httpStatus := builtResponse(builtErrorBodyMsg(err), http.StatusMethodNotAllowed)
	wrapResponse(writer,data,httpStatus)

}



func WrapOkEmptyResponse(writer http.ResponseWriter) {
	wrapResponse(writer,[]byte("{}"),http.StatusOK)
}


func hashPassword(pass string) string {
	hasher := sha256.New()
	hasher.Write([]byte(pass))
	encrypterPassword := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return encrypterPassword
}


