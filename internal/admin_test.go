package internal

import (
	"errors"
	"fmt"
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"golang.org/x/crypto/bcrypt"
	"reflect"
	"testing"
)

func TestProvider_CreateUser(t *testing.T) {
	tests := []struct {
		name                   string
		givenUser              User
		bcryptPasswordError    error
		bcryptPasswordPassword []byte
		dbExpectedUser         storage.User // password not encrypted
		dbReturnError          error
		expectedError          error
	}{
		{
			name: "Happycase",
			givenUser: User{
				EMail:    "test@test.test",
				Password: "s3cr3t",
				Claims:   map[string]interface{}{"cLaIM": "as"},
			},
			dbExpectedUser: storage.User{
				EMail:    "test@test.test",
				Password: []byte("s3cr3t"),
				Claims:   map[string]interface{}{"cLaIM": "as"},
			},
		}, {
			name: "user already exists",
			givenUser: User{
				EMail:    "test@test.test",
				Password: "s3cr3t",
			},
			dbReturnError: storage.ErrUserAlreadyExists,
			dbExpectedUser: storage.User{
				EMail:    "test@test.test",
				Password: []byte("s3cr3t"),
			},
			expectedError: ErrUserAlreadyExists,
		}, {
			name: "Some db error",
			givenUser: User{
				EMail:    "test@test.test",
				Password: "s3cr3t",
			},
			dbReturnError: errors.New("my custom error. ALARM"),
			dbExpectedUser: storage.User{
				EMail:    "test@test.test",
				Password: []byte("s3cr3t"),
			},
			expectedError: errors.New(`failed to query user with email "test@test.test": my custom error. ALARM`),
		}, {
			name: "failed to bcrypt password",
			givenUser: User{
				EMail:    "test@test.test",
				Password: "s3cr3t",
			},
			bcryptPasswordError: errors.New("failed to bcrypt password"),
			dbReturnError:       errors.New("my custom error. ALARM"),
			expectedError:       errors.New(`failed to bcrypt password: failed to bcrypt password`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.bcryptPasswordError != nil {
				oldBcryptPassword := bcryptPassword
				defer func() { bcryptPassword = oldBcryptPassword }()
				bcryptPassword = func(password string) ([]byte, error) {
					return tt.bcryptPasswordPassword, tt.bcryptPasswordError
				}
			}
			var givenDbUser storage.User
			toTest := Provider{
				Storage: &StorageMock{
					CreateUserFunc: func(user storage.User) error {
						givenDbUser = user
						return tt.dbReturnError
					},
				},
			}

			err := toTest.CreateUser(tt.givenUser)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:%s\nGiven:%s", tt.expectedError, err)
			}

			if givenDbUser.EMail != tt.dbExpectedUser.EMail {
				t.Errorf("Given db user > email is not as expected: \nExpected:%s\nGiven:%s", tt.dbExpectedUser.EMail, givenDbUser.EMail)
			}

			err = bcrypt.CompareHashAndPassword(givenDbUser.Password, tt.dbExpectedUser.Password)
			if err != nil && !reflect.DeepEqual(givenDbUser.Password, tt.dbExpectedUser.Password) {
				t.Errorf("Given db user > password is not as expected: \nExpected:%s\nGiven(bcrypted):%s", tt.dbExpectedUser.Password, givenDbUser.Password)
			}
		})
	}
}

func TestProvider_GetUser(t *testing.T) {
	tests := []struct {
		name            string
		givenEMail      string
		dbExpectedEMail string
		dbReturnUser    storage.User
		dbReturnError   error
		expectedError   error
		expectedUser    User
	}{
		{
			name:            "Happycase",
			dbExpectedEMail: "test@test.test",
			dbReturnUser: storage.User{
				EMail: "test.test@test.test",
				Claims: map[string]interface{}{
					"claaa": "bbb",
				},
				Password: []byte("password"),
			},
			givenEMail: "test@test.test",
			expectedUser: User{
				EMail:    "test.test@test.test",
				Password: "**********",
				Claims: map[string]interface{}{
					"claaa": "bbb",
				},
			},
		}, {
			name:            "user not found",
			givenEMail:      "test@test.test",
			dbExpectedEMail: "test@test.test",
			dbReturnError:   storage.ErrUserNotFound,
			expectedError:   ErrUserNotFound,
		}, {
			name:            "Some db error",
			givenEMail:      "test@test.test",
			dbExpectedEMail: "test@test.test",
			dbReturnError:   errors.New("my custom error. ALARM"),
			expectedError:   errors.New(`failed to find user with email "test@test.test": my custom error. ALARM`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var givenEMail string
			toTest := Provider{
				Storage: &StorageMock{
					UserFunc: func(email string) (storage.User, error) {
						givenEMail = email
						return tt.dbReturnUser, tt.dbReturnError
					},
				},
			}

			user, err := toTest.GetUser(tt.givenEMail)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Fatalf("Processing error is not as expected: \nExpected:%s\nGiven:%s", tt.expectedError, err)
			}

			if !reflect.DeepEqual(user, tt.expectedUser) {
				t.Errorf("Returned user is not as expected. Given:\n%#v\nExpected:\n%#v", user, tt.expectedUser)
			}

			if givenEMail != tt.dbExpectedEMail {
				t.Errorf("Given db email is not as expected: \nExpected:%s\nGiven:%s", tt.dbExpectedEMail, givenEMail)
			}
		})
	}

}

func TestProvider_UpdateUser_Happycase(t *testing.T) {
	dbUserToUpdate := storage.User{
		EMail:    "test.test@test.test",
		Password: []byte("testSecret"),
		Claims: map[string]interface{}{
			"c": "g",
		},
	}
	var dbUpdateUser storage.User
	toTest := Provider{
		Storage: &StorageMock{
			UserFunc: func(email string) (storage.User, error) {
				return dbUserToUpdate, nil
			},
			UpdateUserFunc: func(user storage.User) error {
				dbUpdateUser = user
				return nil
			},
		},
	}

	updatedUser, err := toTest.UpdateUser("test.test@test.test", User{
		Password: "newPassword",
		Claims: map[string]interface{}{
			"d": "w",
		},
	})
	if err != nil {
		t.Fatal("unexpected error", err)
	}

	expectedUpdatedUser := User{
		EMail:    "test.test@test.test",
		Password: "**********",
		Claims: map[string]interface{}{
			"d": "w",
		},
	}
	if !reflect.DeepEqual(updatedUser, expectedUpdatedUser) {
		t.Errorf("returned updated user is not as expected. Expected:\n%#v\nGiven:\n%#v", expectedUpdatedUser, updatedUser)
	}

	expectedDBUpdateUser := storage.User{
		EMail:    "test.test@test.test",
		Password: []byte("newPassword"),
		Claims: map[string]interface{}{
			"d": "w",
		},
	}
	if dbUpdateUser.EMail != expectedDBUpdateUser.EMail {
		t.Errorf("user.email to update in db is not as expected. Expected:\n%q\nGiven:\n%q", expectedDBUpdateUser.EMail, dbUpdateUser.EMail)
	}

	err = bcrypt.CompareHashAndPassword(dbUpdateUser.Password, expectedDBUpdateUser.Password)
	if err != nil {
		t.Errorf("user.password to update in db is not as expected. Expected:\n%q", expectedDBUpdateUser.Password)
	}

	if !reflect.DeepEqual(dbUpdateUser.Claims, expectedDBUpdateUser.Claims) {
		t.Errorf("user.claims to update in db is not as expected. Expected:\n%#v\nGiven:\n%#v", expectedDBUpdateUser.Claims, dbUpdateUser.Claims)
	}
}

func TestProvider_UpdateUser_UnableToGetUser(t *testing.T) {
	tests := []struct {
		name                string
		dbUserResponseError error
		expectedError       error
	}{
		{
			name:                "user not found",
			dbUserResponseError: storage.ErrUserNotFound,
			expectedError:       ErrUserNotFound,
		},
		{
			name:                "unexpected error",
			dbUserResponseError: errors.New("nope"),
			expectedError:       errors.New("failed to find user to update: nope"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userEMail := "test.test@test.test"
			var dbUserCalledEMail string
			toTest := Provider{
				Storage: &StorageMock{
					UserFunc: func(email string) (storage.User, error) {
						dbUserCalledEMail = email
						return storage.User{}, tt.dbUserResponseError
					},
				},
			}

			_, err := toTest.UpdateUser(userEMail, User{
				Password: "newPassword",
				Claims: map[string]interface{}{
					"d": "w",
				},
			})
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Errorf("unexpected error. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			}
			if dbUserCalledEMail != userEMail {
				t.Errorf("db called with unexpected email. Expected: %q, Given: %q", userEMail, dbUserCalledEMail)
			}
		})
	}

}

func TestProvider_UpdateUser_UnableToUpdateUser(t *testing.T) {
	tests := []struct {
		name                      string
		dbUpdateUserResponseError error
		expectedError             error
	}{
		{
			name:                      "user not found",
			dbUpdateUserResponseError: storage.ErrUserNotFound,
			expectedError:             ErrUserNotFound,
		},
		{
			name:                      "unexpected error",
			dbUpdateUserResponseError: errors.New("nope"),
			expectedError:             errors.New("failed to update user: nope"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userEMail := "test.test@test.test"
			toTest := Provider{
				Storage: &StorageMock{
					UserFunc: func(_ string) (storage.User, error) {
						return storage.User{
							EMail:    userEMail,
							Password: []byte("bycryptedPassword"),
							Claims: map[string]interface{}{
								"stored": "claim",
							},
						}, nil
					},
					UpdateUserFunc: func(_ storage.User) error {
						return tt.dbUpdateUserResponseError
					},
				},
			}

			_, err := toTest.UpdateUser(userEMail, User{
				Password: "newPassword",
				Claims: map[string]interface{}{
					"d": "w",
				},
			})
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Errorf("unexpected error. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			}
		})
	}
}

func TestProvider_UpdateUser_UnableToBcryptPassword(t *testing.T) {
	oldBcryptPassword := bcryptPassword
	defer func() { bcryptPassword = oldBcryptPassword }()
	bcryptPassword = func(password string) ([]byte, error) {
		return nil, errors.New("failed to bcryptPassword")
	}

	userEMail := "test.test@test.test"
	toTest := Provider{
		Storage: &StorageMock{
			UserFunc: func(_ string) (storage.User, error) {
				return storage.User{
					EMail:    userEMail,
					Password: []byte("bycryptedPassword"),
					Claims: map[string]interface{}{
						"stored": "claim",
					},
				}, nil
			},
			UpdateUserFunc: func(_ storage.User) error {
				return storage.ErrUserNotFound
			},
		},
	}

	_, err := toTest.UpdateUser(userEMail, User{
		Password: "newPassword",
		Claims: map[string]interface{}{
			"d": "w",
		},
	})

	expectedErr := errors.New("failed to bcrypt new password: failed to bcryptPassword")
	if fmt.Sprint(err) != fmt.Sprint(expectedErr) {
		t.Errorf("unexpected error. Expected:\n%q\nGiven:\n%q", expectedErr, err)
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
			name:            "user not found",
			givenEMail:      "test@test.test",
			dbExpectedEMail: "test@test.test",
			dbReturnError:   storage.ErrUserNotFound,
			expectedError:   ErrUserNotFound,
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
