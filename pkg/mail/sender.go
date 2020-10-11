package mail

import (
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
		msg: "Subject: Write Your Subject\n\n" + "This is your Email Body",
	}
}

func (s *Sender) SendEmail(toEmail string) {

	status := smtp.SendMail("smtp.gmail.com:587", smtp.PlainAuth("", s.fromEmail, s.password, "smtp.gmail.com"), s.fromEmail, []string{toEmail}, []byte(s.msg))
	if status != nil {
		log.Printf("Error from SMTP Server: %s", status)
		return
	}
	log.Print("Email Sent Successfully")
}
