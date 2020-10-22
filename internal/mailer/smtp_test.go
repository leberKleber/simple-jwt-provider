package mailer

import (
	"errors"
	"fmt"
	"gopkg.in/mail.v2"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name                    string
		dialerDialSendCloser    mail.SendCloser
		dialerDialErr           error
		loadTemplatesTemplate   mailTemplate
		loadTemplatesErr        error
		expectedErr             error
		expectedMailerTemplates map[string]template
	}{
		{
			name: "Happycase",
			dialerDialSendCloser: &sendCloserMock{
				CloseFunc: func() error { return nil },
			},
			loadTemplatesTemplate: mailTemplate{
				name: "password-reset-request",
			},
			expectedMailerTemplates: map[string]template{
				"password-reset-request": mailTemplate{
					name: "password-reset-request",
				},
			},
		}, {
			name:          "Unable to connect to smtp server",
			dialerDialErr: errors.New("unable to dial: !42"),
			expectedErr:   errors.New("failed to connect to smtp server: unable to dial: !42"),
		}, {
			name: "Unable to load templates",
			dialerDialSendCloser: &sendCloserMock{
				CloseFunc: func() error { return nil },
			},
			loadTemplatesErr: errors.New("angry file system: you're stupid peace of s*it"),
			expectedErr:      errors.New("failed to load password-reset mailTemplate: angry file system: you're stupid peace of s*it"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldBuildDialer := buildDialer
			oldLoadTemplates := loadTemplates
			defer func() {
				buildDialer = oldBuildDialer
				loadTemplates = oldLoadTemplates
			}()

			givenTemplatesFolderPath := "/my/mailTemplate/path"
			givenUsername := ">username<"
			givenPassword := ">password<"
			givenHost := ">host<"
			givenPort := 5555
			givenTLSInsecureSkipVerify := true
			givenTLSServerName := ">tlsServerName<"
			givenDialer := &dialerMock{
				DialFunc: func() (mail.SendCloser, error) {
					return tt.dialerDialSendCloser, tt.dialerDialErr
				},
			}

			buildDialer = func(username, password, host string, port int, tlsInsecureSkipVerify bool, tlsServerName string) dialer {
				return givenDialer
			}

			loadTemplates = func(path, name string) (mailTemplate, error) {
				if path != givenTemplatesFolderPath {
					t.Errorf("unexpected loadTemplates.path. Given: %q, Expected: %q", path, givenTemplatesFolderPath)
				}

				expectedLoadedTemplateName := "password-reset-request"
				if name != expectedLoadedTemplateName {
					t.Errorf("unexpected loadTemplates.name. Given: %q, Expected: %q", name, expectedLoadedTemplateName)
				}

				return tt.loadTemplatesTemplate, tt.loadTemplatesErr
			}

			mailer, err := New(givenTemplatesFolderPath, givenUsername, givenPassword, givenHost, givenPort, givenTLSInsecureSkipVerify, givenTLSServerName)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("Unexpected error. Given:\n%q\nExpected:\n%q", err, tt.loadTemplatesErr)
			} else if err != nil {
				return
			}

			if !reflect.DeepEqual(mailer.templates, tt.expectedMailerTemplates) {
				t.Fatalf("mailer.templates are not as expected. Given:\n%#v\nExpected:\n%#v", mailer.templates, tt.expectedMailerTemplates)
			}

			if !reflect.DeepEqual(mailer.dialer, givenDialer) {
				t.Fatalf("mailer.dialer is not as expected. Given:\n%#v\nExpected:\n%#v", mailer.dialer, givenDialer)
			}
		})
	}

}

func TestMailer_SendPasswordResetRequestEMail_HappyCase(t *testing.T) {
	givenRecipient := ">recipient<"
	givenPasswordResetToken := ">passwordResetToken<"
	givenClaims := map[string]interface{}{
		"customClaim4711": 3,
	}

	prrMail := mail.NewMessage(mail.SetCharset("UTF-8"))
	prrMail.SetHeader("test_id", "yay")

	var mailsToSend []*mail.Message

	dialer := &dialerMock{
		DialAndSendFunc: func(msgs ...*mail.Message) error {
			mailsToSend = msgs
			return nil
		},
	}

	var calledMailData interface{}
	tplID := "password-reset-request"
	tplMock := &templateMock{
		RenderFunc: func(mailData interface{}) (*mail.Message, error) {
			calledMailData = mailData
			return prrMail, nil
		},
	}
	tpls := map[string]template{
		tplID: tplMock,
	}

	m := Mailer{
		dialer:    dialer,
		templates: tpls,
	}

	err := m.SendPasswordResetRequestEMail(givenRecipient, givenPasswordResetToken, givenClaims)
	if err != nil {
		t.Fatal("Unexpected error", err)
	}

	expectedSendMails := []*mail.Message{prrMail}
	if !reflect.DeepEqual(mailsToSend, expectedSendMails) {
		t.Errorf("The send mail(s) are not the rendered. Rendered: %#v. Send: %#v", mailsToSend, expectedSendMails)
	}

	tplMockRenderCalls := tplMock.RenderCalls()
	if len(tplMockRenderCalls) != 1 {
		t.Errorf("tpls[%q].Render should be called 1 time but was %d", tplID, len(tplMockRenderCalls))
	}

	dialerDialAndSendCalls := dialer.DialAndSendCalls()
	if len(dialerDialAndSendCalls) != 1 {
		t.Errorf("dialer.DialAndSendCalls should be called 1 time but was %d", len(dialerDialAndSendCalls))
	}

	expectedMailData := struct {
		Recipient          string
		PasswordResetToken string
		Claims             map[string]interface{}
	}{
		Recipient:          givenRecipient,
		PasswordResetToken: givenPasswordResetToken,
		Claims:             givenClaims,
	}
	if !reflect.DeepEqual(expectedMailData, calledMailData) {
		t.Errorf("called mail data are not as expected. Expected:\n%#v\nGiven:\n%#v", expectedMailData, calledMailData)
	}

	dialerDialCalls := dialer.DialCalls()
	if len(dialerDialCalls) != 0 {
		t.Errorf("dialer.DialAndSendCalls should be called 0 time but was %d", len(dialerDialCalls))
	}
}

func TestMailer_SendPasswordResetRequestEMail_TemplateNotFound(t *testing.T) {
	givenRecipient := ">recipient<"
	givenPasswordResetToken := ">passwordResetToken<"
	givenClaims := map[string]interface{}{
		"customClaim4711": 3,
	}

	prrMail := mail.NewMessage(mail.SetCharset("UTF-8"))
	prrMail.SetHeader("test_id", "yay")

	m := Mailer{
		templates: map[string]template{},
	}

	err := m.SendPasswordResetRequestEMail(givenRecipient, givenPasswordResetToken, givenClaims)
	expectedError := errors.New("could not found mailTemplate with name \"password-reset-request\"")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("Unexpected error. Error:\n%q,\nExpected:\n%q", err, expectedError)
	}
}

func TestMailer_SendPasswordResetRequestEMail_FailedToRenderTemplate(t *testing.T) {
	givenRecipient := ">recipient<"
	givenPasswordResetToken := ">passwordResetToken<"
	givenClaims := map[string]interface{}{
		"customClaim4711": 3,
	}

	prrMail := mail.NewMessage(mail.SetCharset("UTF-8"))
	prrMail.SetHeader("test_id", "yay")

	var calledMailData interface{}
	tplID := "password-reset-request"
	tplMock := &templateMock{
		RenderFunc: func(mailData interface{}) (*mail.Message, error) {
			calledMailData = mailData
			return prrMail, errors.New("i dont think so")
		},
	}
	tpls := map[string]template{
		tplID: tplMock,
	}

	m := Mailer{
		templates: tpls,
	}

	err := m.SendPasswordResetRequestEMail(givenRecipient, givenPasswordResetToken, givenClaims)
	expectedError := errors.New("failed to render mail mailTemplate: i dont think so")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("Unexpected error. Error:\n%q,\nExpected:\n%q", err, expectedError)
	}

	tplMockRenderCalls := tplMock.RenderCalls()
	if len(tplMockRenderCalls) != 1 {
		t.Errorf("tpls[%q].Render should be called 1 time but was %d", tplID, len(tplMockRenderCalls))
	}

	expectedMailData := struct {
		Recipient          string
		PasswordResetToken string
		Claims             map[string]interface{}
	}{
		Recipient:          givenRecipient,
		PasswordResetToken: givenPasswordResetToken,
		Claims:             givenClaims,
	}
	if !reflect.DeepEqual(expectedMailData, calledMailData) {
		t.Errorf("called mail data are not as expected. Expected:\n%#v\nGiven:\n%#v", expectedMailData, calledMailData)
	}
}

func TestMailer_SendPasswordResetRequestEMail_FailedToSendMail(t *testing.T) {
	givenRecipient := ">recipient<"
	givenPasswordResetToken := ">passwordResetToken<"
	givenClaims := map[string]interface{}{
		"customClaim4711": 3,
	}

	prrMail := mail.NewMessage(mail.SetCharset("UTF-8"))
	prrMail.SetHeader("test_id", "yay")

	var mailsToSend []*mail.Message

	dialer := &dialerMock{
		DialAndSendFunc: func(msgs ...*mail.Message) error {
			mailsToSend = msgs
			return errors.New("perhaps yes but no")
		},
	}

	var calledMailData interface{}
	tplID := "password-reset-request"
	tplMock := &templateMock{
		RenderFunc: func(mailData interface{}) (*mail.Message, error) {
			calledMailData = mailData
			return prrMail, nil
		},
	}
	tpls := map[string]template{
		tplID: tplMock,
	}

	m := Mailer{
		dialer:    dialer,
		templates: tpls,
	}

	err := m.SendPasswordResetRequestEMail(givenRecipient, givenPasswordResetToken, givenClaims)
	expectedError := errors.New("failed to send email: perhaps yes but no")
	if fmt.Sprint(err) != fmt.Sprint(expectedError) {
		t.Fatalf("Unexpected error. Error:\n%q,\nExpected:\n%q", err, expectedError)
	}

	expectedSendMails := []*mail.Message{prrMail}
	if !reflect.DeepEqual(mailsToSend, expectedSendMails) {
		t.Errorf("The send mail(s) are not the rendered. Rendered: %#v. Send: %#v", mailsToSend, expectedSendMails)
	}

	tplMockRenderCalls := tplMock.RenderCalls()
	if len(tplMockRenderCalls) != 1 {
		t.Errorf("tpls[%q].Render should be called 1 time but was %d", tplID, len(tplMockRenderCalls))
	}

	dialerDialAndSendCalls := dialer.DialAndSendCalls()
	if len(dialerDialAndSendCalls) != 1 {
		t.Errorf("dialer.DialAndSendCalls should be called 1 time but was %d", len(dialerDialAndSendCalls))
	}

	expectedMailData := struct {
		Recipient          string
		PasswordResetToken string
		Claims             map[string]interface{}
	}{
		Recipient:          givenRecipient,
		PasswordResetToken: givenPasswordResetToken,
		Claims:             givenClaims,
	}
	if !reflect.DeepEqual(expectedMailData, calledMailData) {
		t.Errorf("called mail data are not as expected. Expected:\n%#v\nGiven:\n%#v", expectedMailData, calledMailData)
	}

	dialerDialCalls := dialer.DialCalls()
	if len(dialerDialCalls) != 0 {
		t.Errorf("dialer.DialAndSendCalls should be called 0 time but was %d", len(dialerDialCalls))
	}
}
