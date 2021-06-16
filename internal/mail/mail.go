package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
	"refactory/notes/internal/config"
)

func Sent(name, email string, code int) error {
	dialer := gomail.NewDialer(config.Cfg().MailHost, config.Cfg().MailPort, config.Cfg().MailUser, config.Cfg().MailPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.Cfg().MailUser)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Register Verification")
	mailer.SetBody("text/html", fmt.Sprintf("hello %s, \n this is your verification code: \n %d", name, code))

	if err := dialer.DialAndSend(mailer); nil != err {
		return errors.Wrap(err, fmt.Sprintf("[mailer] Add - Sending email to %s", email))
	}

	return nil
}
