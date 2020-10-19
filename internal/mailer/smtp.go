package mailer

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/mail.v2"
)

//go:generate moq -out send_closer_moq_test.go . sendCloser
type sendCloser mail.SendCloser

//go:generate moq -out dialer_moq_test.go . dialer
type dialer interface {
	DialAndSend(msgs ...*mail.Message) error
	Dial() (mail.SendCloser, error)
}

//go:generate moq -out template_moq_test.go . template
type template interface {
	Render(args interface{}) (*mail.Message, error)
}

type Mailer struct {
	dialer    dialer
	templates map[string]template
}

var buildDialer = func(username string, password string, host string, port int, tlsInsecureSkipVerify bool, tlsServerName string) dialer {
	d := mail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify, ServerName: tlsServerName}

	return d
}

func New(templatesFolderPath, username, password, host string, port int, tlsInsecureSkipVerify bool, tlsServerName string) (*Mailer, error) {
	d := buildDialer(username, password, host, port, tlsInsecureSkipVerify, tlsServerName)

	//check connection and auth
	sc, err := d.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to smtp server: %w", err)
	}
	defer func() { _ = sc.Close() }()

	pwRestTmpl, err := loadTemplates(templatesFolderPath, PasswordResetRequestTemplateName)
	if err != nil {
		return nil, fmt.Errorf("failed to load password-reset mailTemplate: %w", err)
	}

	return &Mailer{
		dialer: d,
		templates: map[string]template{
			PasswordResetRequestTemplateName: pwRestTmpl,
		},
	}, nil
}

func (m *Mailer) SendPasswordResetRequestEMail(recipient, passwordResetToken string, claims map[string]interface{}) error {
	mailData := struct {
		Recipient          string
		PasswordResetToken string
		Claims             map[string]interface{}
	}{
		Recipient:          recipient,
		PasswordResetToken: passwordResetToken,
		Claims:             claims,
	}

	tpl, found := m.templates[PasswordResetRequestTemplateName]
	if !found {
		return fmt.Errorf("could not found mailTemplate with name %q", PasswordResetRequestTemplateName)
	}

	msg, err := tpl.Render(mailData)
	if err != nil {
		return fmt.Errorf("failed to render mail mailTemplate: %w", err)
	}

	err = m.dialer.DialAndSend(msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %s", err)
	}

	return nil
}
