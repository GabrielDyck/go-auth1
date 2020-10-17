package profile

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/backend/auth"
	"auth1/pkg/routes/internal"
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
	authService auth.AuthService
}

func NewProfileInfoService(db mysql.Account, authService auth.AuthService) profileInfoService {
	return profileInfoService{
		db:          db,
		authService: authService,
	}
}

func (s *profileInfoService) editProfileInfo(request *http.Request) ([]byte, int) {
	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}

	token := request.Header.Get("Authorization")

	isAuthorized, err := s.authService.IsProfileEditorAuthorized(int64(id), token)

	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}

	if !isAuthorized {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("you are not authorized to edit a profile that is not belonging to you")), http.StatusBadRequest)
	}

	var req ProfileReq
	err = internal.ParseRequest(request, &req)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}
	log.Println(id)
	log.Println(req)

	account, err := s.db.GetAccountById(int64(id))
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}
	if account == nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}

	if req.Email != account.Email {

		if account.AccountType == api.Google {
			return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("cannot edit email of an Google Account")), http.StatusBadRequest)
		}

		alreadyExist, err := s.db.AccountAlreadyExists(req.Email, account.AccountType)
		if err != nil {
			return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
		}

		if alreadyExist {
			return internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("the email is already in use")), http.StatusBadRequest)
		}

	}

	err = s.db.EditProfileInfo(int64(id), req.Email, req.Address, req.Fullname, req.Phone)
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
	}

	return internal.BuiltResponse(req, http.StatusOK)
}

func (s *profileInfoService) AddRoutes(router *mux.Router) {
	router.HandleFunc(editProfilePath, func(writer http.ResponseWriter, request *http.Request) {
		resp, status := s.editProfileInfo(request)
		writer.WriteHeader(status)
		writer.Write(resp)
	}).Methods("POST")

	router.HandleFunc(profilePath, func(writer http.ResponseWriter, request *http.Request) {
		resp, status := s.getProfileInfo(request)
		writer.WriteHeader(status)
		writer.Write(resp)
	}).Methods("GET")

}

func (s *profileInfoService) getProfileInfo( request *http.Request)([]byte, int) {

	id, err := strconv.Atoi(mux.Vars(request)["id"])
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}
	log.Println(id)

	account, err := s.db.GetAccountById(int64(id))
	if err != nil {
		return internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
	}

	return internal.BuiltResponse(account, http.StatusOK)

}

