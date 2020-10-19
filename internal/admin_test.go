package internal

import (
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestProvider_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		givenEMail     string
		givenPassword  string
		givenClaims    map[string]interface{}
		expectedError  error
		dbExpectedUser storage.User //password not encrypted
		dbReturnError  error
	}{
		{
			name:          "Happycase",
			givenEMail:    "test@test.test",
			givenPassword: "s3cr3t",
			givenClaims:   map[string]interface{}{"cLaIM": "as"},
			dbExpectedUser: storage.User{
				EMail:    "test@test.test",
				Password: []byte("s3cr3t"),
				Claims:   map[string]interface{}{"cLaIM": "as"},
			},
		}, {
			name:          "User already exists",
			givenEMail:    "test@test.test",
			givenPassword: "s3cr3t",
			dbReturnError: storage.ErrUserAlreadyExists,
			dbExpectedUser: storage.User{
				EMail:    "test@test.test",
				Password: []byte("s3cr3t"),
			},
			expectedError: ErrUserAlreadyExists,
		}, {
			name:          "Some db error",
			givenEMail:    "test@test.test",
			givenPassword: "s3cr3t",
			dbReturnError: errors.New("my custom error. ALARM"),
			dbExpectedUser: storage.User{
				EMail:    "test@test.test",
				Password: []byte("s3cr3t"),
			},
			expectedError: errors.New(`failed to query user with email "test@test.test": my custom error. ALARM`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenDbUser storage.User
			toTest := Provider{
				Storage: &StorageMock{
					CreateUserFunc: func(user storage.User) error {
						givenDbUser = user
						return tt.dbReturnError
					},
				},
			}

			err := toTest.CreateUser(tt.givenEMail, tt.givenPassword, tt.givenClaims)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:%s\nGiven:%s", tt.expectedError, err)
			}

			if givenDbUser.EMail != tt.dbExpectedUser.EMail {
				t.Errorf("Given db user > email is not as expected: \nExpected:%s\nGiven:%s", tt.dbExpectedUser.EMail, givenDbUser.EMail)
			}

			if err := bcrypt.CompareHashAndPassword(givenDbUser.Password, tt.dbExpectedUser.Password); err != nil {
				t.Errorf("Given db user > password is not as expected: \nExpected:%s\nGiven(bcrypted):%s", tt.dbExpectedUser.Password, givenDbUser.Password)
			}
		})
	}

}

func TestProvider_DeleteUser(t *testing.T) {
	tests := []struct {
		name            string
		givenEMail      string
		expectedError   error
		dbExpectedEMail string
		dbReturnError   error
	}{
		{
			name:            "Happycase",
			dbExpectedEMail: "test@test.test",
			givenEMail:      "test@test.test",
		}, {
			name:            "User not found",
			givenEMail:      "test@test.test",
			dbExpectedEMail: "test@test.test",
			dbReturnError:   storage.ErrUserNotFound,
			expectedError:   ErrUserNotFound,
		}, {
			name:            "User still has tokens",
			givenEMail:      "test@test.test",
			dbExpectedEMail: "test@test.test",
			dbReturnError:   storage.ErrUserStillHasTokens,
			expectedError:   ErrUserStillHasTokens,
		}, {
			name:            "Some db error",
			givenEMail:      "test@test.test",
			dbExpectedEMail: "test@test.test",
			dbReturnError:   errors.New("my custom error. ALARM"),
			expectedError:   errors.New(`failed to delete user with email "test@test.test": my custom error. ALARM`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenEMail string
			toTest := Provider{
				Storage: &StorageMock{
					DeleteUserFunc: func(email string) error {
						givenEMail = email
						return tt.dbReturnError
					},
				},
			}

			err := toTest.DeleteUser(tt.givenEMail)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:%s\nGiven:%s", tt.expectedError, err)
			}

			if givenEMail != tt.dbExpectedEMail {
				t.Errorf("Given db email is not as expected: \nExpected:%s\nGiven:%s", tt.dbExpectedEMail, givenEMail)
			}
		})
	}

}
