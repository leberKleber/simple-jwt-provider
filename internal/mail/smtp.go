package mail

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/mail.v2"
)

//go:generate moq -out dialer_moq_test.go . dialer
type dialer interface {
	DialAndSend(msgs ...*mail.Message) error
	Dial() (mail.SendCloser, error)
}

type Mailer struct {
	dialer    dialer
	templates map[string]Template
}

func New(templatesFolderPath, username, password, host string, port int, tlsInsecureSkipVerify bool, tlsServerName string) (*Mailer, error) {
	d := buildDialer(username, password, host, port, tlsInsecureSkipVerify, tlsServerName)

	//check connection and auth
	rd, err := d.Dial()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to smtp server: %w", err)
	}
	defer func() { _ = rd.Close() }()

	pwRestTmpl, err := load(templatesFolderPath, PasswordResetRequestTemplateName)
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

func buildDialer(username string, password string, host string, port int, tlsInsecureSkipVerify bool, tlsServerName string) dialer {
	d := mail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: tlsInsecureSkipVerify, ServerName: tlsServerName}

	return d
}

func (m *Mailer) SendPasswordResetRequestEMail(recipient, passwordResetToken string, claims map[string]interface{}) error {
	mailData := struct {
		EMail              string
		PasswordResetToken string
		Claims             map[string]interface{}
	}{
		EMail:              recipient,
		PasswordResetToken: passwordResetToken,
		Claims:             claims,
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
