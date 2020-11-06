package mailer

import (
	"bytes"
	"fmt"
	"gopkg.in/mail.v2"
	"gopkg.in/yaml.v2"
	htmlTemplate "html/template"
	"path/filepath"
	textTemplate "text/template"
)

const passwordResetRequestTemplateName = "password-reset-request"

var htmlTemplateParseFiles = htmlTemplate.ParseFiles
var textTemplateParseFiles = textTemplate.ParseFiles
var ymlTemplateParseFiles = textTemplate.ParseFiles

type mailTemplate struct {
	name       string
	htmlTmpl   *htmlTemplate.Template
	textTmpl   *textTemplate.Template
	headerTmpl *textTemplate.Template
}

var loadTemplates = func(path, name string) (mailTemplate, error) {
	htmlTmpl, err := htmlTemplateParseFiles(filepath.Join(path, fmt.Sprintf("%s.html", name)))
	if err != nil {
		return mailTemplate{}, fmt.Errorf("failed to load mail-body-html mailTemplate: %w", err)
	}

	textTmpl, err := textTemplateParseFiles(filepath.Join(path, fmt.Sprintf("%s.txt", name)))
	if err != nil {
		return mailTemplate{}, fmt.Errorf("failed to load mail-body-text mailTemplate: %w", err)
	}

	headerTmpl, err := ymlTemplateParseFiles(filepath.Join(path, fmt.Sprintf("%s.yml", name)))
	if err != nil {
		return mailTemplate{}, fmt.Errorf("failed to load mail-headers-yml mailTemplate: %w", err)
	}

	return mailTemplate{
		name:       name,
		htmlTmpl:   htmlTmpl,
		textTmpl:   textTmpl,
		headerTmpl: headerTmpl,
	}, nil
}

func (t mailTemplate) Render(args interface{}) (*mail.Message, error) {
	msg := mail.NewMessage()

	err := renderHeaders(msg, t.headerTmpl, args)
	if err != nil {
		return nil, fmt.Errorf("failed to render mail-headers-yml: %w", err)
	}

	var buf bytes.Buffer
	err = t.textTmpl.Execute(&buf, args)
	if err != nil {
		return nil, fmt.Errorf("failed to render mail-body-text: %w", err)
	}
	msg.SetBody("text/plain", buf.String())

	buf.Reset()
	err = t.htmlTmpl.Execute(&buf, args)
	if err != nil {
		return nil, fmt.Errorf("failed to render mail-body-html: %w", err)
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
