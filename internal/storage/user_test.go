package storage

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"testing"
)

func TestStorage_User(t *testing.T) {
	tests := []struct {
		name           string
		givenEMail     string
		dbResponseRows *sqlmock.Rows
		dbResponseErr  error
		expectedUser   User
		expectedError  error
	}{
		{
			name:       "Happycase",
			givenEMail: "info@leberkleber.io",
			dbResponseRows: sqlmock.NewRows([]string{"password", "claims"}).
				AddRow("bcryptedPassword", `{"customClaim1": 4711}`),
			expectedUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
		},
		{
			name:          "No results",
			givenEMail:    "info@leberkleber.io",
			dbResponseErr: sql.ErrNoRows,
			expectedError: ErrUserNotFound,
		},
		{
			name:          "Unexpected db error",
			givenEMail:    "info@leberkleber.io",
			dbResponseErr: errors.New("I used a shitty db"),
			expectedError: errors.New("failed to query user: I used a shitty db"),
		},
		{
			name:       "Non json claims (should not be possible)",
			givenEMail: "info@leberkleber.io",
			dbResponseRows: sqlmock.NewRows([]string{"password", "claims"}).
				AddRow("bcryptedPassword", "customClaim1\n4711}"),
			expectedError: errors.New("failed to unmarshal user>claims: invalid character 'c' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("Failed to create sql mock", err)
			}

			expectedQuery := mock.
				ExpectQuery(`SELECT password, claims FROM users WHERE email = \$1;`).
				WithArgs(tt.givenEMail).
				WillReturnError(tt.dbResponseErr)

			if tt.dbResponseRows != nil {
				expectedQuery.WillReturnRows(tt.dbResponseRows)
			}

			s := Storage{db: db}

			user, err := s.User(tt.givenEMail)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Errorf("Returned error is not as expected. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			}

			//TODO use reflect deep equal
			if fmt.Sprint(user) != fmt.Sprint(tt.expectedUser) {
				t.Errorf("Returned user is not as expected. Expected:\n%#v\nGiven:\n%#v", tt.expectedUser, user)
			}
		})
	}
}

func TestStorage_CreateUser(t *testing.T) {
	tests := []struct {
		name               string
		givenUser          User
		dbResponseErr      error
		expectedDBEMail    string
		expectedDBPassword []byte
		expectedDBClaims   []byte
		expectedError      error
	}{
		{
			name: "Happycase",
			givenUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
			expectedDBEMail:    "info@leberkleber.io",
			expectedDBPassword: []byte("bcryptedPassword"),
			expectedDBClaims:   []byte(`{"customClaim1":4711}`),
		},
		{
			name: "Unexpected db error",
			givenUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
			dbResponseErr:      errors.New("nope"),
			expectedDBEMail:    "info@leberkleber.io",
			expectedDBPassword: []byte("bcryptedPassword"),
			expectedDBClaims:   []byte(`{"customClaim1":4711}`),
			expectedError:      errors.New("failed to exec create stmt: nope"),
		},
		{
			name: "User already exists",
			givenUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
			dbResponseErr: &pq.Error{
				Constraint: "email_unique",
			},
			expectedDBEMail:    "info@leberkleber.io",
			expectedDBPassword: []byte("bcryptedPassword"),
			expectedDBClaims:   []byte(`{"customClaim1":4711}`),
			expectedError:      ErrUserAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("Failed to create sql mock", err)
			}

			mock.
				ExpectExec(`INSERT INTO users \(email, password, claims\) VALUES\(\$1, \$2, \$3\);`).
				WithArgs(tt.expectedDBEMail, tt.expectedDBPassword, tt.expectedDBClaims).
				WillReturnError(tt.dbResponseErr).
				WillReturnResult(sqlmock.NewResult(0, 1))

			s := Storage{db: db}

			err = s.CreateUser(tt.givenUser)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Errorf("Returned error is not as expected. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			}
		})
	}
}

func TestStorage_UpdateUser(t *testing.T) {
	tests := []struct {
		name               string
		givenUser          User
		dbResponseErr      error
		dbResult           driver.Result
		expectedDBEMail    string
		expectedDBPassword []byte
		expectedDBClaims   []byte
		expectedError      error
	}{
		{
			name: "Happycase",
			givenUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
			dbResult:           sqlmock.NewResult(0, 1),
			expectedDBEMail:    "info@leberkleber.io",
			expectedDBPassword: []byte("bcryptedPassword"),
			expectedDBClaims:   []byte(`{"customClaim1":4711}`),
		},
		{
			name: "Unexpected db error",
			givenUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
			dbResult:           sqlmock.NewResult(0, 1),
			dbResponseErr:      errors.New("nope"),
			expectedDBEMail:    "info@leberkleber.io",
			expectedDBPassword: []byte("bcryptedPassword"),
			expectedDBClaims:   []byte(`{"customClaim1":4711}`),
			expectedError:      errors.New("failed to exec update stmt: nope"),
		},
		{
			name: "User already exists",
			givenUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
			dbResult:           sqlmock.NewResult(0, 0),
			expectedDBEMail:    "info@leberkleber.io",
			expectedDBPassword: []byte("bcryptedPassword"),
			expectedDBClaims:   []byte(`{"customClaim1":4711}`),
			expectedError:      ErrUserNotFound,
		},
		{
			name: "Unexpected result error",
			givenUser: User{
				EMail:    "info@leberkleber.io",
				Password: []byte("bcryptedPassword"),
				Claims: map[string]interface{}{
					"customClaim1": 4711,
				}},
			dbResult:           sqlmock.NewErrorResult(errors.New("a random error")),
			expectedDBEMail:    "info@leberkleber.io",
			expectedDBPassword: []byte("bcryptedPassword"),
			expectedDBClaims:   []byte(`{"customClaim1":4711}`),
			expectedError:      errors.New("failed to get count of affected rows: a random error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("Failed to create sql mock", err)
			}

			mock.
				ExpectExec(`UPDATE users SET password = \$2, claims = \$3 WHERE email = \$1;`).
				WithArgs(tt.expectedDBEMail, tt.expectedDBPassword, tt.expectedDBClaims).
				WillReturnError(tt.dbResponseErr).
				WillReturnResult(tt.dbResult)

			s := Storage{db: db}

			err = s.UpdateUser(tt.givenUser)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Errorf("Returned error is not as expected. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			}
		})
	}
}

func TestStorage_DeleteUser(t *testing.T) {
	tests := []struct {
		name                  string
		givenEMail            string
		tokensDBResponseErr   error
		tokensDBResult        driver.Result
		usersDBResponseErr    error
		usersDBResult         driver.Result
		expectedTokensDBEMail string
		expectedUsersDBEMail  string
		expectedError         error
	}{
		{
			name:                  "Happycase",
			givenEMail:            "info@leberkleber.io",
			tokensDBResult:        sqlmock.NewResult(0, 5),
			usersDBResult:         sqlmock.NewResult(0, 1),
			expectedTokensDBEMail: "info@leberkleber.io",
			expectedUsersDBEMail:  "info@leberkleber.io",
		},
		{
			name:                  "Unexpected tokens db error",
			givenEMail:            "info@leberkleber.io",
			tokensDBResponseErr:   errors.New("nope"),
			expectedTokensDBEMail: "info@leberkleber.io",
			expectedError:         errors.New("failed to exec delete tokens from user stmt: nope"),
		},
		{
			name:                  "Unexpected user db error",
			givenEMail:            "info@leberkleber.io",
			tokensDBResult:        sqlmock.NewResult(0, 5),
			usersDBResponseErr:    errors.New("nope"),
			expectedTokensDBEMail: "info@leberkleber.io",
			expectedUsersDBEMail:  "info@leberkleber.io",
			expectedError:         errors.New("failed to exec delete user stmt: nope"),
		},
		{
			name:                  "User doesn't exist",
			givenEMail:            "info@leberkleber.io",
			tokensDBResult:        sqlmock.NewResult(0, 0),
			usersDBResult:         sqlmock.NewResult(0, 0),
			expectedTokensDBEMail: "info@leberkleber.io",
			expectedUsersDBEMail:  "info@leberkleber.io",
			expectedError:         ErrUserNotFound,
		},
		{
			name:                  "Could not get could of affected rows",
			givenEMail:            "info@leberkleber.io",
			tokensDBResult:        sqlmock.NewResult(0, 0),
			usersDBResult:         sqlmock.NewErrorResult(errors.New("a random error")),
			expectedTokensDBEMail: "info@leberkleber.io",
			expectedUsersDBEMail:  "info@leberkleber.io",
			expectedError:         errors.New("failed to get count of affected rows: a random error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("Failed to create sql mock", err)
			}

			mock.ExpectBegin()

			mock.
				ExpectExec(`DELETE FROM tokens WHERE email = \$1;`).
				WithArgs(tt.expectedTokensDBEMail).
				WillReturnError(tt.tokensDBResponseErr).
				WillReturnResult(tt.tokensDBResult)

			mock.
				ExpectExec(`DELETE FROM users WHERE email = \$1;`).
				WithArgs(tt.expectedUsersDBEMail).
				WillReturnError(tt.usersDBResponseErr).
				WillReturnResult(tt.usersDBResult)

			if tt.expectedError == nil {
				mock.ExpectCommit()
			}

			s := Storage{db: db}

			err = s.DeleteUser(tt.givenEMail)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedError) {
				t.Errorf("Returned error is not as expected. Expected:\n%q\nGiven:\n%q", tt.expectedError, err)
			}
		})
	}
}
