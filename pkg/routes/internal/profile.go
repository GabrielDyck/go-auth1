package internal

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

const (
	profilePath          = "/profile-info/{id}"
	editProfilePath          = "/edit-profile/{id}"

)

type ProfileReq struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type profileInfoService struct {
	db mysql.Account
}

func NewProfileInfoService(db mysql.Account) profileInfoService {
	return profileInfoService{db: db}
}

func GetProfileInfo(router *mux.Router, db mysql.Account) {

	service := NewProfileInfoService(db)
	router.HandleFunc(profilePath, func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(mux.Vars(request)["id"])
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		fmt.Println(id)

		account, err := service.getProfileInfo(int64(id))
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		data, httpStatus := builtResponse(account, http.StatusOK)
		wrapResponse(writer, data, httpStatus)
	}).Methods("GET")

}

func EditProfileInfo(router *mux.Router, service profileInfoService) {
	router.HandleFunc(editProfilePath, func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(mux.Vars(request)["id"])
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}

		var req ProfileReq
		err = parseRequest(writer, request, &req)
		if err != nil {
			return
		}
		fmt.Println(id)
		fmt.Println(req)

		err = validateRequiredFields(req)
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		account, err := service.getProfileInfo(int64(id))
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		data, httpStatus := builtResponse(account, http.StatusOK)
		wrapResponse(writer, data, httpStatus)

	}).Methods("POST")

}

func (s *profileInfoService) getProfileInfo(id int64) (*api.Account, error) {
	return s.db.GetAccountById(id)
}

func validateRequiredFields(req ProfileReq) error {

	if req.Email == "" {
		return errors.New("email cannot be empty")
	}
	return nil
}
