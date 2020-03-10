// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package internal

import (
	"sync"
)

var (
	lockMailerMockSendPasswordResetRequestEMail sync.RWMutex
)

// Ensure, that MailerMock does implement Mailer.
// If this is not the case, regenerate this file with moq.
var _ Mailer = &MailerMock{}

// MailerMock is a mock implementation of Mailer.
//
//     func TestSomethingThatUsesMailer(t *testing.T) {
//
//         // make and configure a mocked Mailer
//         mockedMailer := &MailerMock{
//             SendPasswordResetRequestEMailFunc: func(recipient string, passwordResetLink string) error {
// 	               panic("mock out the SendPasswordResetRequestEMail method")
//             },
//         }
//
//         // use mockedMailer in code that requires Mailer
//         // and then make assertions.
//
//     }
type MailerMock struct {
	// SendPasswordResetRequestEMailFunc mocks the SendPasswordResetRequestEMail method.
	SendPasswordResetRequestEMailFunc func(recipient string, passwordResetLink string) error

	// calls tracks calls to the methods.
	calls struct {
		// SendPasswordResetRequestEMail holds details about calls to the SendPasswordResetRequestEMail method.
		SendPasswordResetRequestEMail []struct {
			// Recipient is the recipient argument value.
			Recipient string
			// PasswordResetLink is the passwordResetLink argument value.
			PasswordResetLink string
		}
	}
}

// SendPasswordResetRequestEMail calls SendPasswordResetRequestEMailFunc.
func (mock *MailerMock) SendPasswordResetRequestEMail(recipient string, passwordResetLink string) error {
	if mock.SendPasswordResetRequestEMailFunc == nil {
		panic("MailerMock.SendPasswordResetRequestEMailFunc: method is nil but Mailer.SendPasswordResetRequestEMail was just called")
	}
	callInfo := struct {
		Recipient         string
		PasswordResetLink string
	}{
		Recipient:         recipient,
		PasswordResetLink: passwordResetLink,
	}
	lockMailerMockSendPasswordResetRequestEMail.Lock()
	mock.calls.SendPasswordResetRequestEMail = append(mock.calls.SendPasswordResetRequestEMail, callInfo)
	lockMailerMockSendPasswordResetRequestEMail.Unlock()
	return mock.SendPasswordResetRequestEMailFunc(recipient, passwordResetLink)
}

// SendPasswordResetRequestEMailCalls gets all the calls that were made to SendPasswordResetRequestEMail.
// Check the length with:
//     len(mockedMailer.SendPasswordResetRequestEMailCalls())
func (mock *MailerMock) SendPasswordResetRequestEMailCalls() []struct {
	Recipient         string
	PasswordResetLink string
} {
	var calls []struct {
		Recipient         string
		PasswordResetLink string
	}
	lockMailerMockSendPasswordResetRequestEMail.RLock()
	calls = mock.calls.SendPasswordResetRequestEMail
	lockMailerMockSendPasswordResetRequestEMail.RUnlock()
	return calls
}
