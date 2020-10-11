package mail

import (
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
	return Sender{
		fromEmail: "gabrielapptester@gmail.com",
		password : os.Getenv("SMTP_PASS"),
		msg: "Subject: Auth1 Reset Password\n\n" + "Click here to reset your password http://localhost:9290/resetPassword?token=%s",
	}
}

func (s *Sender) SendEmail(toEmail, token string) {

	status := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", s.fromEmail, s.password, "smtp.gmail.com"), s.fromEmail, []string{toEmail}, []byte(fmt.Sprintf(s.msg,token)))
	if status != nil {
		log.Printf("Error from SMTP Server: %s", status)
		return
	}
	log.Print("Email Sent Successfully")
}
