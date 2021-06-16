package mail

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"github.com/prometheus/common/log"
	"gopkg.in/gomail.v2"
	"refactory/notes/internal/app/model"
	"refactory/notes/internal/app/repository"
	"refactory/notes/internal/config"
)

type Mailer interface {
	Add(session model.Session)
}

type mailer struct {
	repo     repository.UserRepository
	sessions []model.Session
}

func newMessage(email, name string, code int) *gomail.Message {
	message := gomail.NewMessage()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.Cfg().MailUser)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Register Verification")
	mailer.SetBody("text/html", fmt.Sprintf("hello %s, \n this is your verification code: \n %d", name, code))

	return message
}

func NewMailer(repo repository.UserRepository) Mailer {
	mail := &mailer{repo: repo}
	go mail.loop()
	return mail
}

func (c *mailer) Add(session model.Session) {
	c.sessions = append(c.sessions, session)
	log.Infof("added %s to queue mailer", session.Username)
}

func (c *mailer) loop() {
	for {
		if len(c.sessions) > 0 {
			for _, s := range c.sessions {
				log.Infof("sending verification email to %s", s.Email)
				// sent verification email
				message := gomail.NewMessage()
				message.SetHeader("From", config.Cfg().MailUser)
				message.SetHeader("To", s.Email)
				message.SetHeader("Subject", "Register Verification")
				message.SetBody("text/html", fmt.Sprintf("hello %s, \n this is your verification code: \n %d", s.Username, s.Code))

				if err := c.sent(message); nil == err {
					// removing session from queue when success updating session
					s.IsSent = true
					if err := c.repo.UpdateSession(context.Background(), s); nil == err {
						c.sessions = c.sessions[1:]
					} else {
						log.Errorf("error while updating session for %s error %s", s.Email, err.Error())
					}
				} else {
					log.Errorf("error while sending email: %s error: %s", s.Email, err.Error())
				}
			}
		}
	}
}

func (c *mailer) sent(message *gomail.Message) error {
	dialer := gomail.NewDialer(config.Cfg().MailHost, config.Cfg().MailPort, config.Cfg().MailUser, config.Cfg().MailPassword)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := dialer.DialAndSend(message); nil != err {
		return errors.Wrap(err, fmt.Sprintf("[mailer] Add - Sending email to %s", message.GetHeader("To")))
	}

	return nil
}
