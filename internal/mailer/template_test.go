package mailer

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/DusanKasan/parsemail"
	htmlTemplate "html/template"
	"os"
	"reflect"
	"testing"
	textTemplate "text/template"
)

func TestLoadTemplates(t *testing.T) {
	tests := []struct {
		name                           string
		givenTemplatePath              string
		givenTemplateName              string
		parseHTMLFileExpectedFilenames []string
		parseHTMLFileTemplate          *htmlTemplate.Template
		parseHTMLFileErr               error
		parseTextFileExpectedFilenames []string
		parseTextFileTemplate          *textTemplate.Template
		parseTextFileErr               error
		parseYMLFileExpectedFilenames  []string
		parseYMLFileTemplate           *textTemplate.Template
		parseYMLFileErr                error
		expectedTemplate               mailTemplate
		expectedErr                    error
	}{
		{
			name:                           "happycase",
			givenTemplatePath:              "/my/mailTemplate/path",
			givenTemplateName:              "myTemplate",
			parseHTMLFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.html"},
			parseHTMLFileTemplate:          htmlTemplate.New("myTemplateHTML"),
			parseTextFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.txt"},
			parseTextFileTemplate:          textTemplate.New("myTemplateText"),
			parseYMLFileExpectedFilenames:  []string{"/my/mailTemplate/path/myTemplate.yml"},
			parseYMLFileTemplate:           textTemplate.New("myTemplateYML"),
			expectedTemplate: mailTemplate{
				name:       "myTemplate",
				htmlTmpl:   htmlTemplate.New("myTemplateHTML"),
				textTmpl:   textTemplate.New("myTemplateText"),
				headerTmpl: textTemplate.New("myTemplateYML"),
			},
		}, {
			name:                           "could not load html mailTemplate",
			givenTemplatePath:              "/my/mailTemplate/path",
			givenTemplateName:              "myTemplate",
			parseHTMLFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.html"},
			parseHTMLFileErr:               os.ErrNotExist,
			expectedErr:                    errors.New("failed to load mail html body template: file does not exist"),
		}, {
			name:                           "could not load text mailTemplate",
			givenTemplatePath:              "/my/mailTemplate/path",
			givenTemplateName:              "myTemplate",
			parseHTMLFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.html"},
			parseHTMLFileTemplate:          htmlTemplate.New("myTemplateHTML"),
			parseTextFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.txt"},
			parseTextFileErr:               os.ErrPermission,
			expectedErr:                    errors.New("failed to load mail text body template: permission denied"),
		}, {
			name:                           "could not load yml mailTemplate",
			givenTemplatePath:              "/my/mailTemplate/path",
			givenTemplateName:              "myTemplate",
			parseHTMLFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.html"},
			parseHTMLFileTemplate:          htmlTemplate.New("myTemplateHTML"),
			parseTextFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.txt"},
			parseTextFileTemplate:          textTemplate.New("myTemplateText"),
			parseYMLFileExpectedFilenames:  []string{"/my/mailTemplate/path/myTemplate.yml"},
			parseYMLFileErr:                errors.New("abc error"),
			expectedErr:                    errors.New("failed to load mail headers template: abc error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldHTMLTemplateParseFiles := htmlTemplateParseFiles
			oldTextTemplateParseFiles := textTemplateParseFiles
			oldYMLTemplateParseFiles := ymlTemplateParseFiles
			defer func() {
				htmlTemplateParseFiles = oldHTMLTemplateParseFiles
				textTemplateParseFiles = oldTextTemplateParseFiles
				ymlTemplateParseFiles = oldYMLTemplateParseFiles
			}()

			htmlTemplateParseFiles = func(filenames ...string) (*htmlTemplate.Template, error) {
				if !reflect.DeepEqual(filenames, tt.parseHTMLFileExpectedFilenames) {
					t.Errorf("html mailTemplate file path is not as expected. Expected: %q, given: %q", tt.parseHTMLFileExpectedFilenames, filenames)
				}
				return tt.parseHTMLFileTemplate, tt.parseHTMLFileErr
			}
			textTemplateParseFiles = func(filenames ...string) (*textTemplate.Template, error) {
				if !reflect.DeepEqual(filenames, tt.parseTextFileExpectedFilenames) {
					t.Errorf("text mailTemplate file path is not as expected. Expected: %q, given: %q", tt.parseTextFileExpectedFilenames, filenames)
				}
				return tt.parseTextFileTemplate, tt.parseTextFileErr
			}
			ymlTemplateParseFiles = func(filenames ...string) (*textTemplate.Template, error) {
				if !reflect.DeepEqual(filenames, tt.parseYMLFileExpectedFilenames) {
					t.Errorf("yml mailTemplate file path is not as expected. Expected: %q, given: %q", tt.parseYMLFileExpectedFilenames, filenames)
				}
				return tt.parseYMLFileTemplate, tt.parseYMLFileErr
			}

			template, err := loadTemplates(tt.givenTemplatePath, tt.givenTemplateName)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Fatalf("Unexpected error. Given:\n%q\nExpected:\n%q", err, tt.expectedErr)
			} else if err != nil {
				return
			}

			if !reflect.DeepEqual(template, tt.expectedTemplate) {
				t.Fatalf("Unexpected mailTemplate. Given:\n%#v\nExpected:\n%#v", template, tt.expectedTemplate)
			}
		})
	}
}

func TestTemplate_Render(t *testing.T) {
	tests := []struct {
		name             string
		htmlTplContent   string
		textTplContent   string
		headerTplContent string
		givenRenderArgs  interface{}
		expectedError    error
	}{
		{
			name:           "Happycase",
			htmlTplContent: "html template {{.TestID}}",
			textTplContent: "text template {{.TestID}}",
			headerTplContent: `MyHeader:
  - "headaaaaa {{.TestID}}"`,
			givenRenderArgs: struct {
				TestID string
			}{
				"myTestID",
			},
		},
		{
			name: "header-template execute error handling",
			headerTplContent: `MyHeader:
  - "headaaaaa {{.notExisting}}"`,
			givenRenderArgs: struct {
				notExisting string
			}{},
			expectedError: errors.New("failed to render mail headers: template: htmlTemplate:2:17: executing \"htmlTemplate\" at <.notExisting>: notExisting is an unexported field of struct type struct { notExisting string }"),
		},
		{
			name: "invalid header-template yml syntax",
			headerTplContent: `MyHeader:
  "headaaaaa"`,
			expectedError: errors.New("failed to render mail headers: yaml: unmarshal errors:\n  line 2: cannot unmarshal !!str `headaaaaa` into []string"),
		},
		{
			name:           "text-template Execute error handling",
			htmlTplContent: "html template",
			textTplContent: "text template {{.testID}}",
			headerTplContent: `MyHeader:
  - "headaaaaa"`,
			givenRenderArgs: struct {
				testID string
			}{
				"myTestID",
			},
			expectedError: errors.New("failed to render mail text body: template: htmlTemplate:1:16: executing \"htmlTemplate\" at <.testID>: testID is an unexported field of struct type struct { testID string }"),
		},
		{
			name:           "html-template Execute error handling",
			htmlTplContent: "html template {{.testID}}",
			textTplContent: "text template",
			headerTplContent: `MyHeader:
  - "headaaaaa"`,
			givenRenderArgs: struct {
				testID string
			}{
				"myTestID",
			},
			expectedError: errors.New("failed to render mail html body: template: htmlTemplate:1:16: executing \"htmlTemplate\" at <.testID>: testID is an unexported field of struct type struct { testID string }"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlTpl, err := htmlTemplate.New("htmlTemplate").Parse(tt.htmlTplContent)
			if err != nil {
				t.Fatal("failed to parse test html template", err)
			}
			textTpl, err := textTemplate.New("htmlTemplate").Parse(tt.textTplContent)
			if err != nil {
				t.Fatal("failed to parse test html template", err)
			}
			headerTpl, err := textTemplate.New("htmlTemplate").Parse(tt.headerTplContent)
			if err != nil {
				t.Fatal("failed to parse test html template", err)
			}

			mt := mailTemplate{
				name:       "testMail",
				htmlTmpl:   htmlTpl,
				textTmpl:   textTpl,
				headerTmpl: headerTpl,
			}

			mail, err := mt.Render(tt.givenRenderArgs)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("unexpected error while render template. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			} else if err != nil {
				return
			}

			var bb bytes.Buffer
			_, err = mail.WriteTo(&bb)
			if err != nil {
				t.Error("failed to write mail to bb", err)
			}

			parsedEMail, err := parsemail.Parse(&bb) // returns Email struct and error
			if err != nil {
				t.Fatal("failed to parse written mail", err)
			}

			expectedHTMLBody := "html template myTestID"
			if expectedHTMLBody != parsedEMail.HTMLBody {
				t.Errorf("html body is not as expected. Expected: %q, Give: %q", expectedHTMLBody, parsedEMail.HTMLBody)
			}
			expectedTextBody := "text template myTestID"
			if expectedTextBody != parsedEMail.TextBody {
				t.Errorf("text body is not as expected. Expected: %q, Give: %q", expectedTextBody, parsedEMail.TextBody)
			}
			expectedTestHeaderContent := "headaaaaa myTestID"
			testHeaderValue := parsedEMail.Header.Get("MyHeader")
			if expectedTestHeaderContent != testHeaderValue {
				t.Errorf("test header value is not as expected. Expected: %q, Give: %q", expectedTestHeaderContent, testHeaderValue)
			}
		})
	}
}
