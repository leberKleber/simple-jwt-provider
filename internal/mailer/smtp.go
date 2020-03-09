package mailer

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"net/smtp"
)

type Mailer struct {
	dialer    *gomail.Dialer
	templates map[string]Template
}

func New(templateFolderPath, username, password, host string, port int) (*Mailer, error) {
	d := gomail.NewDialer(host, port, username, password)
	d.Auth = &loginAuth{username: username, password: password}

	pwRestTmpl, err := Load(templateFolderPath, PasswordResetRequestTemplateName)
	if err != nil {
		return nil, fmt.Errorf("failed to load password-reset tempate: %w", err)
	}

	return &Mailer{
		dialer: d,
		templates: map[string]Template{
			PasswordResetRequestTemplateName: pwRestTmpl,
		},
	}, nil
}

func (m *Mailer) SendPasswordResetRequestEMail(recipient, passwordResetLink string) error {
	mailData := struct {
		EMail             string
		PasswordResetLink string
	}{
		EMail:             recipient,
		PasswordResetLink: passwordResetLink,
	}

	msg, err := m.templates[PasswordResetRequestTemplateName].Render(mailData)
	if err != nil {
		return fmt.Errorf("failed to render mail template: %w", err)
	}
	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %s", err)
	}

	return nil
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
