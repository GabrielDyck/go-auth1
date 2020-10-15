package internal

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"errors"
	"github.com/gorilla/mux"
	"log"
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

func (s *profileInfoService) GetProfileInfo(router *mux.Router) {

	router.HandleFunc(profilePath, func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(mux.Vars(request)["id"])
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		log.Println(id)

		account, err := s.getProfileInfo(int64(id))
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		data, httpStatus := BuiltResponse(account, http.StatusOK)
		WrapResponse(writer, data, httpStatus)
	}).Methods("GET")

}

func (s *profileInfoService) EditProfileInfo(router *mux.Router) {
	router.HandleFunc(editProfilePath, func(writer http.ResponseWriter, request *http.Request) {
		id, err := strconv.Atoi(mux.Vars(request)["id"])
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}

		var req ProfileReq
		err = ParseRequest(writer, request, &req)
		if err != nil {
			return
		}
		log.Println(id)
		log.Println(req)

		err = validateRequiredFields(req)
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		account, err := s.getProfileInfo(int64(id))
		if err != nil {
			WrapBadRequestResponse(writer, err)
			return
		}
		data, httpStatus := BuiltResponse(account, http.StatusOK)
		WrapResponse(writer, data, httpStatus)

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
