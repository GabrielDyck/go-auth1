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
	profilePath     = "/profile-info/{id}"
	editProfilePath = "/edit-profile/{id}"
)

type ProfileReq struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
}

type profileInfoService struct {
	db          mysql.Account
	authService AuthService
}

func NewProfileInfoService(db mysql.Account, authService AuthService) profileInfoService {
	return profileInfoService{
		db:          db,
		authService: authService,
	}
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

		token:= request.Header.Get("Authorization")

		isAuthorized,err:= s.authService.IsProfileAuthorized(int64(id),token)

		if err!= nil{
			WrapInternalErrorResponse(writer,err)
			return
		}

		if !isAuthorized{
			WrapBadRequestResponse(writer,errors.New("you are not authorized to edit a profile that is not belonging to you"))
			return
		}

		var req ProfileReq
		err = ParseRequest(writer, request, &req)
		if err != nil {
			return
		}
		log.Println(id)
		log.Println(req)

		account, err := s.getProfileInfo(int64(id))
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}
		if account == nil{
			WrapBadRequestResponse(writer, err)
			return
		}

		if req.Email != account.Email {

			if account.AccountType == api.Google{
				WrapBadRequestResponse(writer, errors.New("cannot edit email of an Google Account"))
				return
			}

			alreadyExist, err:= s.db.AccountAlreadyExists(req.Email, account.AccountType)
			if err != nil {
				WrapInternalErrorResponse(writer, err)
				return
			}

			if alreadyExist{
				WrapBadRequestResponse(writer, errors.New("the email "))
				return
			}

		}

		err= s.db.EditProfileInfo(int64(id),req.Email,req.Address,req.Fullname,req.Phone)
		if err != nil {
			WrapInternalErrorResponse(writer, err)
			return
		}

		data, httpStatus := BuiltResponse(req, http.StatusOK)
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
