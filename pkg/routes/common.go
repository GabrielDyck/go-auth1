package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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


func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}