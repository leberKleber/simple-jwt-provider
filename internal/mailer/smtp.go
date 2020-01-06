package mailer

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/smtp"
	"path/filepath"
	"text/template"
)

type Mailer struct {
	d    *gomail.Dialer
	from string
}

func New(fromAddress, username, password, host string, port int) *Mailer {
	d := gomail.NewDialer(host, port, username, password)
	d.Auth = &loginAuth{username: username, password: password}

	return &Mailer{
		d:    d,
		from: fromAddress,
	}
}

func (m *Mailer) SendPasswordResetMail(recipient, passwordResetLink string) error {
	mailData := struct {
		EMail             string
		PasswordResetLink string
	}{
		EMail:             recipient,
		PasswordResetLink: passwordResetLink,
	}

	tpl, err := parseTemplate("password-reset", mailData)
	if err != nil {
		return fmt.Errorf("failed to parse mail template: %w", err)
	}

	msg := newMessage()
	msg.SetFrom(m.from)
	msg.SetTo(recipient)
	msg.SetBodyText(tpl)
	msg.SetSubject("Simple-JWT-Provider: password-reset")

	err = m.d.DialAndSend(msg.Get())
	if err != nil {
		return fmt.Errorf("failed to send email: %s", err)
	}

	return nil
}

func parseTemplate(tplName string, args interface{}) (string, error) {
	txtTplFileName := tplName + ".txt"
	b, err := ioutil.ReadFile(filepath.Join("/home/mmarch/workspace/leberKleber@github.com/simple-jwt-auth/mail-templates", txtTplFileName))
	if err != nil {
		return "", err
	}

	txtTpl, err := template.New(tplName).Parse(string(b))
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = txtTpl.Execute(&buf, args)

	return buf.String(), err
}

type loginAuth struct {
	username, password string
}

func (a *loginAuth) Start(_ *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown response (%s) from server when attempting to use loginAuth", string(fromServer))
		}
	}
	return nil, nil
}
