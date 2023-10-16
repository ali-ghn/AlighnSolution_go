package email

import (
	"crypto/tls"
	"net/smtp"
	"strconv"
	"strings"

	"gopkg.in/gomail.v2"
)

type IEmailSender interface {
	Send(subject string, content string, from string, to string) error
	SendBulk(subject string, content string, from string, to []string) error
}

type EmailSender struct {
	SmtpHost string
	Password string
}

func NewEmailSender(smtpHost string, password string) EmailSender {
	return EmailSender{
		SmtpHost: smtpHost,
		Password: password,
	}
}

func (es EmailSender) Send(subject string, content string, from string, to string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", content)
	smtpHostname := strings.Split(es.SmtpHost, ":")
	port, err := strconv.Atoi(smtpHostname[1])
	if err != nil {
		return err
	}
	hostname := smtpHostname[0]
	n := gomail.NewDialer(hostname, port, from, es.Password)

	n.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err = n.DialAndSend(msg)

	if err != nil {
		return err
	}

	return nil
}

func (es EmailSender) SendBulk(subject string, content string, from string, to []string) error {
	auth := smtp.PlainAuth("", from, es.Password, es.SmtpHost)
	err := smtp.SendMail(es.SmtpHost, auth, from, to, []byte(content))
	return err
}
