package routes

import (
	"auth1/pkg/mysql"
	"auth1/pkg/mysql/model"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
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

type profileInfoService struct {
	db mysql.ProfileInfo
}

func NewProfileInfoService(db mysql.ProfileInfo) profileInfoService {
	return profileInfoService{db: db}
}
func getProfileInfo(router *mux.Router, db mysql.ProfileInfo) {

	service := NewProfileInfoService(db)
	router.HandleFunc(profile, func(writer http.ResponseWriter, request *http.Request) {
		id := mux.Vars(request)["id"]

		fmt.Println(id)

	}).Methods("GET")

	router.HandleFunc(profile, func(writer http.ResponseWriter, request *http.Request) {
		id,err := strconv.Atoi( mux.Vars(request)["id"])
		if err != nil {
			wrapBadRequestResponse(writer, err)
			return
		}

		var req ProfileWriteReq
		err = parseRequest(writer, request, &req)
		if err != nil {
			return
		}
		fmt.Println(id)
		fmt.Println(req)

		err = validateRequiredFields(req)
		if err != nil {
			wrapBadRequestResponse(writer, err)
			return
		}
		account, err :=service.getProfileInfo(int64(id))
		if err != nil {
			wrapBadRequestResponse(writer, err)
			return
		}
		data, httpStatus :=builtResponse(account, http.StatusOK)
		wrapResponse(writer,data,httpStatus)

	}).Methods("POST")

}


func (s *profileInfoService) getProfileInfo(id int64) (*model.Account,error) {
	return s.db.GetProfileInfoById(id)
}


func validateRequiredFields(req ProfileWriteReq) error {

	if req.Email == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}
