package singup

import (
	"auth1/api"
	"auth1/pkg/mysql"
	"auth1/pkg/routes/internal"
	"crypto/rand"
	"fmt"
	"log"
)

type signupService struct {
	db mysql.SignUp
}

func NewSignUpService(db mysql.SignUp) signupService {
	return signupService{
		db: db,
	}
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
func (s *signupService) accountAlreadyExists(email string) (bool, error) {
	log.Println("AccountAlreadyExists")

	return s.db.AccountAlreadyExists(email)
}

func (s *signupService) getProfileInfoByEmailAndAccountType(email string, accountType api.AccountType) (*api.Account, error) {
	log.Println("GetProfileInfoByEmailAndAccountType")
	return s.db.GetProfileInfoByEmailAndAccountType(email, accountType)
}
