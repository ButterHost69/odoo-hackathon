package utils

import (
	"gopkg.in/gomail.v2"
)

var (
	FROM_EMAIL string
	FROM_EMAIL_PASSWORD string
)

func InitEmailClient(fromEmail, fromEmail_password string){
	FROM_EMAIL = fromEmail
	FROM_EMAIL_PASSWORD = fromEmail_password
}

func SMTP_SendMessagetoEmail(email string, subject string,body string) error {
	m := gomail.NewMessage()
    m.SetHeader("From", FROM_EMAIL)
    m.SetHeader("To", email)
    m.SetHeader("Subject", subject)
    m.SetBody("text/plain", body)

    d := gomail.NewDialer("smtp.gmail.com", 587, FROM_EMAIL, FROM_EMAIL_PASSWORD)

    // Send the email
    if err := d.DialAndSend(m); err != nil {
        return err
    }
    return nil
}