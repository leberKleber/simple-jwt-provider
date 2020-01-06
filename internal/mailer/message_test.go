package mailer

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMessage(t *testing.T) {
	message := newMessage()

	assert.NotNil(t, message)
}

func TestMessage_SetFrom(t *testing.T) {
	message := newMessage()

	message.SetFrom("someFrom")

	result := message.Get()

	assert.Equal(t, "someFrom", result.GetHeader("From")[0])
}

func TestMessage_SetTo(t *testing.T) {
	message := newMessage()

	message.SetTo("someTo", "someOtherTo")

	result := message.Get()

	assert.Equal(t, "someTo", result.GetHeader("To")[0])
	assert.Equal(t, "someOtherTo", result.GetHeader("To")[1])
}

func TestMessage_SetCc(t *testing.T) {
	message := newMessage()

	message.SetCc("someCc", "someOtherCc")

	result := message.Get()

	assert.Equal(t, "someCc", result.GetHeader("Cc")[0])
	assert.Equal(t, "someOtherCc", result.GetHeader("Cc")[1])
}

func TestMessage_SetBcc(t *testing.T) {
	message := newMessage()

	message.SetBcc("someBcc", "someOtherBcc")

	result := message.Get()

	assert.Equal(t, "someBcc", result.GetHeader("Bcc")[0])
	assert.Equal(t, "someOtherBcc", result.GetHeader("Bcc")[1])
}

func TestMessage_SetReplyTo(t *testing.T) {
	message := newMessage()

	message.SetReplyTo("someReplyTo")

	result := message.Get()

	assert.Equal(t, "someReplyTo", result.GetHeader("Reply-To")[0])
}

func TestMessage_SetSubject(t *testing.T) {
	message := newMessage()

	message.SetSubject("someSubject")

	result := message.Get()

	assert.Equal(t, "someSubject", result.GetHeader("Subject")[0])
}

func TestMessage_SetBodyText(t *testing.T) {
	message := newMessage()

	message.SetBodyText("someTextBody")

	result := new(bytes.Buffer)
	message.Get().WriteTo(result)

	assert.Contains(t, result.String(), "someTextBody")
}

func TestMessage_SetBodyHtml(t *testing.T) {
	message := newMessage()

	message.SetBodyHtml("<html><body>someHtmlBody</body></html>")

	result := new(bytes.Buffer)
	message.Get().WriteTo(result)

	assert.Contains(t, result.String(), "<html><body>someHtmlBody</body></html>")
}
