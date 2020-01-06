package mailer

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
)

type Message struct {
	m *gomail.Message
}

func newMessage() *Message {
	return &Message{
		m: gomail.NewMessage(),
	}
}

func (m *Message) SetFrom(from string) {
	m.m.SetHeader("From", from)
}

func (m *Message) SetTo(to ...string) {
	m.m.SetHeader("To", to...)
}

func (m *Message) SetCc(cc ...string) {
	m.m.SetHeader("Cc", cc...)
}

func (m *Message) SetBcc(bcc ...string) {
	m.m.SetHeader("Bcc", bcc...)
}

func (m *Message) SetReplyTo(replyTo string) {
	m.m.SetHeader("Reply-To", replyTo)
}

func (m *Message) SetSubject(subject string) {
	m.m.SetHeader("Subject", subject)
}

func (m *Message) SetBodyText(body string) {
	m.m.AddAlternative("text/plain", body)
}

func (m *Message) SetBodyHtml(body string) {
	m.m.AddAlternative("text/html", body)
}

func (m *Message) AddAttachment(fileName string, content []byte) {
	m.m.Attach(
		fileName,
		gomail.SetCopyFunc(func(w io.Writer) error {
			n, err := w.Write(content)
			if err != nil {
				return err
			}
			if n != len(content) {
				return fmt.Errorf("written byte does not match content: got %v, want %v: %w", n, len(content), err)
			}
			return err
		}),
	)
}

func (m *Message) Get() *gomail.Message {
	return m.m
}
