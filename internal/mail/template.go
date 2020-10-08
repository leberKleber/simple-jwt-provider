package mail

import (
	"bytes"
	"fmt"
	"gopkg.in/mail.v2"
	"gopkg.in/yaml.v2"
	htmlTemplate "html/template"
	"path/filepath"
	textTemplate "text/template"
)

var PasswordResetRequestTemplateName = "password-reset-request"

type Template struct {
	name       string
	htmlTmpl   *htmlTemplate.Template
	textTmpl   *textTemplate.Template
	headerTmpl *textTemplate.Template
}

func load(path, name string) (Template, error) {
	htmlTmpl, err := htmlTemplate.ParseFiles(
		filepath.Join(path, fmt.Sprintf("%s.html", name)),
	)
	if err != nil {
		return Template{}, fmt.Errorf("failed to load html template: %w", err)
	}
	textTmpl, err := textTemplate.ParseFiles(
		filepath.Join(path, fmt.Sprintf("%s.txt", name)),
	)
	if err != nil {
		return Template{}, fmt.Errorf("failed to load txt template: %w", err)
	}
	headerTmpl, err := textTemplate.ParseFiles(
		filepath.Join(path, fmt.Sprintf("%s.yml", name)),
	)
	if err != nil {
		return Template{}, fmt.Errorf("failed to load header template: %w", err)
	}

	return Template{
		name:       name,
		htmlTmpl:   htmlTmpl,
		textTmpl:   textTmpl,
		headerTmpl: headerTmpl,
	}, nil
}

func (t Template) Render(args interface{}) (*mail.Message, error) {
	msg := mail.NewMessage()

	err := renderHeaders(msg, t.headerTmpl, args)
	if err != nil {
		return nil, fmt.Errorf("failed to render headers: %w", err)
	}

	var buf bytes.Buffer
	err = t.textTmpl.Execute(&buf, args)
	if err != nil {
		return nil, fmt.Errorf("failed to render mail-text-body: %w", err)
	}
	msg.AddAlternative("text/plain", buf.String())

	buf.Reset()
	err = t.htmlTmpl.Execute(&buf, args)
	if err != nil {
		return nil, fmt.Errorf("failed to render mail-html-body: %w", err)
	}
	msg.AddAlternative("text/html", buf.String())

	return msg, nil
}

func renderHeaders(msg *mail.Message, template *textTemplate.Template, args interface{}) error {
	var buf bytes.Buffer
	err := template.Execute(&buf, args)
	if err != nil {
		return err
	}

	headers := make(map[string][]string)
	err = yaml.NewDecoder(&buf).Decode(&headers)
	if err != nil {
		return err
	}
	msg.SetHeaders(headers)
	return nil
}
