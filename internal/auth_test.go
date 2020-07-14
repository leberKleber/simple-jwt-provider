package internal

import (
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"reflect"
	"testing"
)

func TestProvider_Login(t *testing.T) {
	tests := []struct {
		name                   string
		givenEMail             string
		givenPassword          string
		expectedError          error
		expectedJWT            string
		generatorExpectedEMail string
		generatorJWT           string
		generatorError         error
		dbReturnError          error
		dbReturnUser           storage.User
	}{
		{
			name:                   "Happycase",
			givenEMail:             "test@test.test",
			givenPassword:          "password",
			generatorExpectedEMail: "test@test.test",
			generatorJWT:           "myJWT",
			expectedJWT:            "myJWT",
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
			expectedError: errors.New("failed to query user with email \"not@existing.user\": unexpected error"),
			dbReturnError: errors.New("unexpected error"),
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
			var givenGeneratorEMail string
			var givenGeneratorUserClaims map[string]interface{}
			toTest := Provider{
				Storage: &StorageMock{
					UserFunc: func(email string) (storage.User, error) {
						givenStorageEMail = email
						return tt.dbReturnUser, tt.dbReturnError
					},
				},
				JWTGenerator: &JWTGeneratorMock{
					GenerateFunc: func(email string, userClaims map[string]interface{}) (string, error) {
						givenGeneratorEMail = email
						givenGeneratorUserClaims = userClaims
						return tt.generatorJWT, tt.generatorError
					},
				},
			}

			jwt, err := toTest.Login(tt.givenEMail, tt.givenPassword)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:\n%s\nGiven:\n%s", tt.expectedError, err)
			}

			if jwt != tt.expectedJWT {
				t.Errorf("Given jwt is not as expected: \nExpected:%s\nGiven:%s", tt.expectedJWT, jwt)
			}

			if givenStorageEMail != tt.givenEMail {
				t.Errorf("DB-Requestest User>Email ist not as expected: \nExpected:%s\nGiven:%s", tt.givenEMail, givenStorageEMail)
			}

			if givenGeneratorEMail != tt.generatorExpectedEMail {
				t.Errorf("Generator.Generate email ist not as expected: \nExpected:%s\nGiven:%s", tt.givenEMail, givenGeneratorEMail)
			}

			if !reflect.DeepEqual(givenGeneratorUserClaims, tt.dbReturnUser.Claims) {
				t.Errorf("Generator.Generate userClaims are not as expected: \nExpected:\n%#v\nGiven:\n%#v", tt.givenEMail, givenGeneratorEMail)
			}
		})
	}

}
