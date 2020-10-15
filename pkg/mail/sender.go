package mail

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type Sender struct {
	fromEmail string
	password  string
	msg       string
}

func NewSender() Sender {

	domain:= os.Getenv("DOMAIN")
	return Sender{
		fromEmail: "gabrielapptester@gmail.com",
		password : os.Getenv("SMTP_PASS"),
		msg: "Subject: Auth1 Reset Password\n\n" + "Click here to reset your password http://" + domain +"/reset-password?token=%s",
	}
}

func (s *Sender) SendEmail(toEmail, token string)error {

	status := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", s.fromEmail, s.password, "smtp.gmail.com"), s.fromEmail, []string{toEmail}, []byte(fmt.Sprintf(s.msg,token)))
	if status != nil {
		log.Printf("Error from SMTP Server: %s", status)
		return errors.New("coulnd't send email for recover password")
	}
	log.Print("Email Sent Successfully")
	return nil
}
