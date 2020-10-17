package signup

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/oauth"
	"auth1/pkg/routes/internal"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type signupService struct {
	db mysql.SignUp
}

func NewSignUpService(db mysql.SignUp) signupService {
	return signupService{
		db: db,
	}
}


func (s * signupService) signUp( request *http.Request) ([]byte,int,string){

	var req api.UserSignReq
	err := internal.ParseRequest(request, &req)
	if err != nil {
		data,status:=internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
		return data,status,""

	}
	log.Println(req)

	var account *api.Account

	switch req.AccountType {

	case api.Basic:
		account, err = s.getProfileInfoByEmailAndAccountType(req.Email, api.Basic)

		if err != nil {
			data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			return data,status,""

		}

		if account!=nil {
			data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("user already registered")), http.StatusBadRequest)
			return data,status,""

		}

		err = s.signUpBasicAccount(req)
		if err != nil {
			data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			return data,status,""

		}

		account, err = s.getProfileInfoByEmailAndAccountType(req.Email, api.Basic)
		if err != nil {
			data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			return data,status,""

		}

	case api.Google:
		tokenInfo,err := oauth.VerifyIdToken(req.GoogleToken)
		if err != nil {
			data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusBadRequest)
			return data,status,""

		}

		account, err = s.getProfileInfoByEmailAndAccountType(tokenInfo.Email, api.Google)
		if err != nil {
			data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			return data,status,""

		}

		if account == nil {
			err=s.signUpGoogleAccount(tokenInfo.Email)
			if err != nil {
				data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
				return data,status,""

			}
		}

		account, err = s.getProfileInfoByEmailAndAccountType(tokenInfo.Email, api.Google)
		if err != nil {
			data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
			return data,status,""

		}
	default:
		data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(errors.New("unknown account type")), http.StatusBadRequest)
		return data,status,""

	}

	token, err := s.generateSessionToken(account.ID)
	if err != nil {
		data,status:= internal.BuiltResponse(internal.BuiltErrorBodyMsg(err), http.StatusInternalServerError)
		return data,status,""

	}

	data,status:=internal.BuiltResponse(account, http.StatusOK)
	return  data,status,token
}

func (s *signupService) signUpBasicAccount(req api.UserSignReq) error {
	log.Println("signUpBasicAccount")

	hashedPassword := internal.HashPassword(req.Password)
	return s.db.SignUpBasicAccount(req.Email, hashedPassword)
}

// TODO : extract with signin
func (s *signupService) generateSessionToken(id int64) (string, error) {
	log.Println("generateSessionToken")

	token := make([]byte, 128)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenString := fmt.Sprintf("%X",token)
	return s.createAuthToken(id, tokenString)
}

func (s *signupService) createAuthToken(id int64, tokenString string) (string, error) {
	err := s.db.CreateAuthorizationToken(id, tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *signupService) signUpGoogleAccount(email string) error {
	log.Println("SignUpGoogleAccount")
	return s.db.SignUpGoogleAccount(email)
}
func (s *signupService) accountAlreadyExists(email string, accountType api.AccountType) (bool, error) {
	log.Println("AccountAlreadyExists")

	return s.db.AccountAlreadyExists(email,accountType)
}

func (s *signupService) getProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error) {
	log.Println("GetProfileInfoByEmailAndAccountType")
	return s.db.GetProfileInfoByEmailAndAccountType(email, accountType)
}
