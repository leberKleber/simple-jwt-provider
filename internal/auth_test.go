package internal

import (
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"regexp"
	"testing"
	"time"
)

func TestProvider_Login(t *testing.T) {
	bcryptCost = bcrypt.MinCost

	tests := []struct {
		name                      string
		givenEMail                string
		givenPassword             string
		expectedError             error
		expectedAccessToken       string
		expectedRefreshToken      string
		generatorExpectedEMail    string
		generateAccessToken       string
		generateAccessTokenError  error
		generateRefreshToken      string
		generateRefreshTokenError error
		dbReturnError             error
		dbReturnUser              storage.User
	}{
		{
			name:                   "Happycase",
			givenEMail:             "test@test.test",
			givenPassword:          "password",
			generatorExpectedEMail: "test@test.test",
			generateAccessToken:    "myJWT",
			generateRefreshToken:   "myRefreshJWT",
			expectedAccessToken:    "myJWT",
			expectedRefreshToken:   "myRefreshJWT",
			dbReturnUser: storage.User{
				Password: []byte("$2a$12$1v7O.pNLqugJjcePyxvUj.GK37YoAbJvSW/9bULSRmq5C4SkoU2OO"),
				EMail:    "test@test.test",
				Claims: map[string]interface{}{
					"myCustomClaim": "value",
				},
			},
		},
		{
			name:          "User not found",
			givenEMail:    "not@existing.user",
			givenPassword: "password",
			expectedError: ErrUserNotFound,
			dbReturnError: storage.ErrUserNotFound,
		},
		{
			name:          "Unexpected db error",
			givenEMail:    "not@existing.user",
			givenPassword: "password",
			expectedError: errors.New("failed to find user with email \"not@existing.user\": unexpected error"),
			dbReturnError: errors.New("unexpected error"),
		},
		{
			name: "Failed to generate accessToken",
			dbReturnUser: storage.User{
				Password: []byte("$2a$12$1v7O.pNLqugJjcePyxvUj.GK37YoAbJvSW/9bULSRmq5C4SkoU2OO"),
				EMail:    "test@test.test",
				Claims: map[string]interface{}{
					"myCustomClaim": "value",
				},
			},
			givenEMail:               "not@existing.user",
			givenPassword:            "password",
			generateAccessTokenError: errors.New("error 42"),
			expectedError:            errors.New("failed to generate access-token: error 42"),
		},
		{
			name: "Failed to generate refreshToken",
			dbReturnUser: storage.User{
				Password: []byte("$2a$12$1v7O.pNLqugJjcePyxvUj.GK37YoAbJvSW/9bULSRmq5C4SkoU2OO"),
				EMail:    "test@test.test",
				Claims: map[string]interface{}{
					"myCustomClaim": "value",
				},
			},
			givenEMail:                "not@existing.user",
			givenPassword:             "password",
			generateRefreshTokenError: errors.New("error 42"),
			expectedError:             errors.New("failed to generate refresh-token: error 42"),
		},
		{
			name:          "Incorrect Password",
			givenEMail:    "test@test.test",
			givenPassword: "wrongPassword",
			expectedError: ErrIncorrectPassword,
			dbReturnUser: storage.User{
				Password: []byte("$2a$12$1v7O.pNLqugJjcePyxvUj.GK37YoAbJvSW/9bULSRmq5C4SkoU2OO"),
				EMail:    "test@test.test",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenStorageEMail string
			var givenGenerateRefreshTokenEMail string
			var givenGenerateAccessTokenEMail string
			var givenGenerateAccessTokenUserClaims map[string]interface{}
			toTest := Provider{
				Storage: &StorageMock{
					UserFunc: func(email string) (storage.User, error) {
						givenStorageEMail = email
						return tt.dbReturnUser, tt.dbReturnError
					},
				},
				JWTProvider: &JWTProviderMock{
					GenerateAccessTokenFunc: func(email string, userClaims map[string]interface{}) (string, error) {
						givenGenerateAccessTokenEMail = email
						givenGenerateAccessTokenUserClaims = userClaims
						return tt.generateAccessToken, tt.generateAccessTokenError
					},
					GenerateRefreshTokenFunc: func(email string) (string, error) {
						givenGenerateRefreshTokenEMail = email
						return tt.generateRefreshToken, tt.generateRefreshTokenError
					},
				},
			}

			accessToken, refreshToken, err := toTest.Login(tt.givenEMail, tt.givenPassword)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:\n%s\nGiven:\n%s", tt.expectedError, err)
			} else if err != nil {
				return
			}

			if accessToken != tt.expectedAccessToken {
				t.Errorf("Given accessToken is not as expected: \nExpected:%s\nGiven:%s", tt.expectedAccessToken, accessToken)
			}

			if refreshToken != tt.expectedRefreshToken {
				t.Errorf("Given refreshToken is not as expected: \nExpected:%s\nGiven:%s", tt.expectedRefreshToken, refreshToken)
			}

			if givenStorageEMail != tt.givenEMail {
				t.Errorf("DB-Requestest User>Email ist not as expected: \nExpected:%s\nGiven:%s", tt.givenEMail, givenStorageEMail)
			}

			if givenGenerateAccessTokenEMail != tt.generatorExpectedEMail {
				t.Errorf("Generator.GenerateAccessToken email ist not as expected: \nExpected:%s\nGiven:%s", tt.givenEMail, givenGenerateAccessTokenEMail)
			}

			if !reflect.DeepEqual(givenGenerateAccessTokenUserClaims, tt.dbReturnUser.Claims) {
				t.Errorf("Generator.GenerateAccessToken userClaims are not as expected: \nExpected:\n%#v\nGiven:\n%#v", tt.generatorExpectedEMail, givenGenerateAccessTokenEMail)
			}

			if givenGenerateRefreshTokenEMail != tt.generatorExpectedEMail {
				t.Errorf("Generator.GenerateRefreshToken email ist not as expected: \nExpected:%s\nGiven:%s", tt.generatorExpectedEMail, givenGenerateRefreshTokenEMail)
			}
		})
	}

}

func TestProvider_CreatePasswordResetRequest(t *testing.T) {
	tests := []struct {
		name                      string
		givenEMail                string
		expectedError             error
		dbUserReturnError         error
		dbCreateTokenReturnError  error
		generateHEXTokenError     error
		mailerError               error
		dbExpectedToken           storage.Token
		expectedMailRecipient     string
		passwordResetTokenPresent bool
	}{
		{
			name:                  "Happycase",
			givenEMail:            "test.test@test.test",
			expectedMailRecipient: "test.test@test.test",
			dbExpectedToken: storage.Token{
				Type:  "reset",
				EMail: "test.test@test.test",
				ID:    0,
			},
			expectedError: nil,
		}, {
			name:              "User not found",
			givenEMail:        "not@existing.user",
			dbUserReturnError: storage.ErrUserNotFound,
			expectedError:     ErrUserNotFound,
		}, {
			name:              "Unexpected db error while finding user",
			givenEMail:        "test.test@test",
			dbUserReturnError: errors.New("random error"),
			expectedError:     errors.New("failed to find user with email \"test.test@test\": random error"),
		}, {
			name:                     "Unexpected db error while create token",
			givenEMail:               "test.test@test",
			dbCreateTokenReturnError: errors.New("random error"),
			expectedError:            errors.New("failed to create password reset token for email \"test.test@test\": random error"),
			dbExpectedToken: storage.Token{
				Type:  "reset",
				EMail: "test.test@test",
				ID:    0,
			},
		}, {
			name:                  "Mailer error",
			givenEMail:            "test.test@test",
			mailerError:           errors.New("random error"),
			expectedError:         errors.New("failed to send password reset email: random error"),
			expectedMailRecipient: "test.test@test",
			dbExpectedToken: storage.Token{
				Type:  "reset",
				EMail: "test.test@test",
				ID:    0,
			},
		}, {
			name:                  "Unable to generate HEX token",
			givenEMail:            "test.test@test",
			generateHEXTokenError: errors.New("random error"),
			expectedError:         errors.New("failed to generate password reset token: random error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.generateHEXTokenError != nil {
				oldGenerateHEXToken := generateHEXToken
				defer func() { generateHEXToken = oldGenerateHEXToken }()

				generateHEXToken = func() (string, error) {
					return "", tt.generateHEXTokenError
				}
			}

			var storageUserEMail string
			var storageCreateTokenToken storage.Token
			var mailerRecipient string
			var mailerPasswordResetToken string
			toTest := Provider{
				Storage: &StorageMock{
					UserFunc: func(email string) (storage.User, error) {
						storageUserEMail = email
						return storage.User{}, tt.dbUserReturnError
					},
					CreateTokenFunc: func(t storage.Token) (int64, error) {
						storageCreateTokenToken = t
						return 0, tt.dbCreateTokenReturnError
					},
				},
				Mailer: &MailerMock{
					SendPasswordResetRequestEMailFunc: func(recipient string, passwordResetToken string, claims map[string]interface{}) error {
						mailerRecipient = recipient
						mailerPasswordResetToken = passwordResetToken
						return tt.mailerError
					},
				},
			}

			err := toTest.CreatePasswordResetRequest(tt.givenEMail)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:\n%s\nGiven:\n%s", tt.expectedError, err)
			}

			if storageUserEMail != tt.givenEMail {
				t.Errorf("The sorage requested usermail is not as expected: \nExpected:\n%s\nGiven:\n%s", tt.givenEMail, storageUserEMail)
			}

			storageCreateTokenToken.Token = ""
			storageCreateTokenToken.CreatedAt = time.Time{}
			if !reflect.DeepEqual(storageCreateTokenToken, tt.dbExpectedToken) {
				t.Errorf("The sorage token to create is not as expected: \nExpected:\n%#v\nGiven:\n%#v", tt.dbExpectedToken, storageCreateTokenToken)
			}

			if mailerRecipient != tt.expectedMailRecipient {
				t.Errorf("The mailer recipient is not as expected: \nExpected:\n%#v\nGiven:\n%#v", tt.expectedMailRecipient, mailerRecipient)
			}

			if tt.passwordResetTokenPresent {
				matched, err := regexp.Match("^[0-9A-Fa-f]{64}$", []byte(mailerPasswordResetToken))
				if err != nil {
					t.Fatalf("could not compile regex")
				}
				if !matched {
					t.Errorf("PasswordResetToken should be a 64 char hex string but was %q", mailerPasswordResetToken)
				}
			}
		})
	}
}

func TestProvider_ResetPassword(t *testing.T) {
	bcryptCost = bcrypt.MinCost

	tests := []struct {
		name                string
		givenEMail          string
		givenResetToken     string
		givenNewPassword    string
		bcryptPasswordError error
		dbToken             []storage.Token
		dbTokenError        error
		dbUser              storage.User
		dbUserError         error
		dbUpdateUserError   error
		dbDeleteTokenError  error
		expectedError       error
	}{
		{
			name:             "Happycase",
			givenNewPassword: "newPassword",
			givenResetToken:  "resetToken",
			givenEMail:       "email",
			dbToken: []storage.Token{
				{ID: 4, CreatedAt: time.Now(), Token: "myToken1", Type: "reset", EMail: "email"},
				{ID: 5, CreatedAt: time.Now(), Token: "myToken2", Type: "other", EMail: "email"},
			},
		},
		{
			name:             "No token found",
			givenNewPassword: "newPassword",
			givenResetToken:  "resetToken",
			givenEMail:       "email",
			expectedError:    ErrNoValidTokenFound,
			dbToken:          []storage.Token{},
		},
		{
			name:             "Error while find tokens",
			givenNewPassword: "newPassword",
			givenResetToken:  "resetToken",
			givenEMail:       "email",
			expectedError:    errors.New("failed to find tokens: unexpected error"),
			dbToken:          []storage.Token{},
			dbTokenError:     errors.New("unexpected error"),
		},
		{
			name:             "Error while find user",
			givenNewPassword: "newPassword",
			givenResetToken:  "resetToken",
			givenEMail:       "email",
			dbToken: []storage.Token{
				{ID: 4, CreatedAt: time.Now(), Token: "myToken1", Type: "reset", EMail: "email"},
				{ID: 5, CreatedAt: time.Now(), Token: "myToken2", Type: "other", EMail: "email"},
			},
			dbUserError:   errors.New("unexpected error"),
			expectedError: errors.New("failed to find user with email \"email\": unexpected error"),
		},
		{
			name:             "Error while update user",
			givenNewPassword: "newPassword",
			givenResetToken:  "resetToken",
			givenEMail:       "email",
			dbToken: []storage.Token{
				{ID: 4, CreatedAt: time.Now(), Token: "myToken1", Type: "reset", EMail: "email"},
				{ID: 5, CreatedAt: time.Now(), Token: "myToken2", Type: "other", EMail: "email"},
			},
			dbUpdateUserError: errors.New("unexpected error"),
			expectedError:     errors.New("failed to update user: unexpected error"),
		},
		{
			name:             "Error while delete token",
			givenNewPassword: "newPassword",
			givenResetToken:  "resetToken",
			givenEMail:       "email",
			dbToken: []storage.Token{
				{ID: 4, CreatedAt: time.Now(), Token: "myToken1", Type: "reset", EMail: "email"},
				{ID: 5, CreatedAt: time.Now(), Token: "myToken2", Type: "other", EMail: "email"},
			},
			dbDeleteTokenError: errors.New("unexpected error"),
			expectedError:      errors.New("failed to delete token: unexpected error"),
		},
		{
			name:                "Error bcrypt password",
			givenNewPassword:    "newPassword",
			givenResetToken:     "resetToken",
			givenEMail:          "email",
			bcryptPasswordError: errors.New("something went wrong"),
			dbToken: []storage.Token{
				{ID: 4, CreatedAt: time.Now(), Token: "myToken1", Type: "reset", EMail: "email"},
			},
			expectedError: errors.New("failed to bcrypt password: something went wrong"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.bcryptPasswordError != nil {
				oldBcryptPassword := bcryptPassword
				defer func() { bcryptPassword = oldBcryptPassword }()

				bcryptPassword = func(password string) ([]byte, error) {
					return nil, tt.bcryptPasswordError
				}
			}

			toTest := Provider{
				Storage: &StorageMock{
					TokensByEMailAndTokenFunc: func(email string, token string) ([]storage.Token, error) {
						return tt.dbToken, tt.dbTokenError
					},
					UserFunc: func(email string) (storage.User, error) {
						return tt.dbUser, tt.dbUserError
					},
					UpdateUserFunc: func(user storage.User) error {
						return tt.dbUpdateUserError
					},
					DeleteTokenFunc: func(id int64) error {
						return tt.dbDeleteTokenError
					},
				},
			}

			err := toTest.ResetPassword(tt.givenEMail, tt.givenResetToken, tt.givenNewPassword)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:\n%s\nGiven:\n%s", tt.expectedError, err)
			}
		})
	}
}
