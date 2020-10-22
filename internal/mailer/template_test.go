package mailer

import (
	"errors"
	"fmt"
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
			expectedErr:                    errors.New("failed to load mail-body-html mailTemplate: file does not exist"),
		}, {
			name:                           "could not load text mailTemplate",
			givenTemplatePath:              "/my/mailTemplate/path",
			givenTemplateName:              "myTemplate",
			parseHTMLFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.html"},
			parseHTMLFileTemplate:          htmlTemplate.New("myTemplateHTML"),
			parseTextFileExpectedFilenames: []string{"/my/mailTemplate/path/myTemplate.txt"},
			parseTextFileErr:               os.ErrPermission,
			expectedErr:                    errors.New("failed to load mail-body-text mailTemplate: permission denied"),
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
			expectedErr:                    errors.New("failed to load mail-headers-yml mailTemplate: abc error"),
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
					t.Errorf("text mailTemplate file path is not as expected. Expected: %q, given: %q", tt.parseHTMLFileExpectedFilenames, filenames)
				}
				return tt.parseTextFileTemplate, tt.parseTextFileErr
			}
			ymlTemplateParseFiles = func(filenames ...string) (*textTemplate.Template, error) {
				if !reflect.DeepEqual(filenames, tt.parseYMLFileExpectedFilenames) {
					t.Errorf("yml mailTemplate file path is not as expected. Expected: %q, given: %q", tt.parseHTMLFileExpectedFilenames, filenames)
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

}
