package storage

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"testing"
	"time"
)

func TestStorage_CreateToken(t *testing.T) {
	tests := []struct {
		name                string
		givenToken          Token
		dbResponseErr       error
		dbResponseRows      *sqlmock.Rows
		expectedDBEMail     string
		expectedDBToken     string
		expectedDBType      string
		expectedDBCreatedAt time.Time
		expectedID          int64
		expectedErr         error
	}{
		{
			name: "Happycase",
			givenToken: Token{
				EMail:     "info@leberkleber.io",
				CreatedAt: time.Date(2020, 2, 1, 4, 46, 45, 2, time.UTC),
				Token:     "myGeneratedToken",
				Type:      "reset",
			},
			dbResponseRows:      sqlmock.NewRows([]string{"id"}).AddRow(42),
			expectedDBEMail:     "info@leberkleber.io",
			expectedDBType:      "reset",
			expectedDBToken:     "myGeneratedToken",
			expectedDBCreatedAt: time.Date(2020, 2, 1, 4, 46, 45, 2, time.UTC),
			expectedID:          42,
		},
		{
			name: "Unexpected db error",
			givenToken: Token{
				EMail:     "info@leberkleber.io",
				CreatedAt: time.Date(2020, 2, 1, 4, 46, 45, 2, time.UTC),
				Token:     "myGeneratedToken",
				Type:      "reset",
			},
			dbResponseRows:      sqlmock.NewRows([]string{"id"}).AddRow(42),
			dbResponseErr:       errors.New("nope"),
			expectedDBEMail:     "info@leberkleber.io",
			expectedDBType:      "reset",
			expectedDBToken:     "myGeneratedToken",
			expectedDBCreatedAt: time.Date(2020, 2, 1, 4, 46, 45, 2, time.UTC),
			expectedErr:         errors.New("failed to exec create token stmt: nope"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("Failed to create sql mock", err)
			}

			expectedQuery := mock.
				ExpectQuery(`INSERT INTO tokens \(email, token, type, created_at\) VALUES\(\$1, \$2, \$3, \$4\) RETURNING id;`).
				WithArgs(tt.expectedDBEMail, tt.expectedDBToken, tt.expectedDBType, tt.expectedDBCreatedAt).
				WillReturnError(tt.dbResponseErr)
			if tt.dbResponseRows != nil {
				expectedQuery.WillReturnRows(tt.dbResponseRows)
			}

			s := Storage{db: db}

			id, err := s.CreateToken(tt.givenToken)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Errorf("Returned error is not as expected. Expected:\n%q\nGiven:\n%q", tt.expectedErr, err)
			}
			if id != tt.expectedID {
				t.Errorf("Returned id is not as expected. Expected:\n%d\nGiven:\n%d", tt.expectedID, id)
			}
		})
	}
}

func TestStorage_TokensByEMailAndToken(t *testing.T) {
	tests := []struct {
		name            string
		givenEMail      string
		givenToken      string
		dbResponseErr   error
		dbResponseRows  *sqlmock.Rows
		expectedDBEMail string
		expectedDBToken string
		expectedTokens  []Token
		expectedErr     error
	}{
		{
			name:       "Happycase",
			givenEMail: "info@leberkleber.io",
			givenToken: "myGeneratedToken",
			dbResponseRows: sqlmock.NewRows([]string{"id", "type", "created_at"}).
				AddRow(1, "reset", time.Date(2020, 01, 01, 01, 01, 01, 01, time.UTC)).
				AddRow(42, "reset", time.Date(1999, 01, 01, 01, 01, 01, 01, time.UTC)),
			expectedDBEMail: "info@leberkleber.io",
			expectedDBToken: "myGeneratedToken",
			expectedTokens: []Token{
				{ID: 1, EMail: "info@leberkleber.io", Type: "reset", CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 01, time.UTC), Token: "myGeneratedToken"},
				{ID: 42, EMail: "info@leberkleber.io", Type: "reset", CreatedAt: time.Date(1999, 01, 01, 01, 01, 01, 01, time.UTC), Token: "myGeneratedToken"},
			},
		},
		{
			name:            "Error while exec stmt",
			givenEMail:      "info@leberkleber.io",
			givenToken:      "myGeneratedToken",
			dbResponseErr:   errors.New("nope"),
			expectedDBEMail: "info@leberkleber.io",
			expectedDBToken: "myGeneratedToken",
			expectedErr:     errors.New("failed to exec select token stmt: nope"),
		},
		{
			name:       "Unable to scan sql response",
			givenEMail: "info@leberkleber.io",
			givenToken: "myGeneratedToken",
			dbResponseRows: sqlmock.NewRows([]string{"id", "type"}).
				AddRow(1, "reset").
				AddRow(42, "reset"),
			expectedDBEMail: "info@leberkleber.io",
			expectedDBToken: "myGeneratedToken",
			expectedErr:     errors.New("failed to scan select token stmt result: sql: expected 2 destination arguments in Scan, not 3"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("Failed to create sql mock", err)
			}

			expectedQuery := mock.
				ExpectQuery(`SELECT id, type, created_at FROM tokens WHERE email = \$1 AND token = \$2;`).
				WithArgs(tt.expectedDBEMail, tt.expectedDBToken).
				WillReturnError(tt.dbResponseErr)
			if tt.dbResponseRows != nil {
				expectedQuery.WillReturnRows(tt.dbResponseRows)
			}

			s := Storage{db: db}

			tokens, err := s.TokensByEMailAndToken(tt.givenEMail, tt.givenToken)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Errorf("Returned error is not as expected. Expected:\n%q\nGiven:\n%q", tt.expectedErr, err)
			}
			if !reflect.DeepEqual(tokens, tt.expectedTokens) {
				t.Errorf("Returned id is not as expected. Expected:\n%#v\nGiven:\n%#v", tt.expectedTokens, tokens)
			}
		})
	}
}

func TestStorage_DeleteToken(t *testing.T) {
	tests := []struct {
		name             string
		givenID          int64
		dbResponseErr    error
		dbResponseResult driver.Result
		expectedDBID     int64
		expectedErr      error
	}{
		{
			name:             "Happycase",
			givenID:          5561,
			dbResponseResult: sqlmock.NewResult(0, 1),
			expectedDBID:     5561,
		},
		{
			name:          "Error while exec",
			givenID:       5561,
			expectedDBID:  5561,
			dbResponseErr: errors.New("nope"),
			expectedErr:   errors.New("failed to delete token: nope"),
		},
		{
			name:             "Error while get affected rows (should not be possible)",
			givenID:          5561,
			expectedDBID:     5561,
			dbResponseResult: sqlmock.NewErrorResult(errors.New("aaaaaaaaaa")),
			expectedErr:      errors.New("could not get num of affected row: aaaaaaaaaa"),
		},
		{
			name:             "Token not found",
			givenID:          5561,
			expectedDBID:     5561,
			dbResponseResult: sqlmock.NewResult(0, 0),
			expectedErr:      ErrTokenNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal("Failed to create sql mock", err)
			}

			mock.
				ExpectExec(`DELETE FROM tokens WHERE id = \$1`).
				WithArgs(tt.expectedDBID).
				WillReturnResult(tt.dbResponseResult).
				WillReturnError(tt.dbResponseErr)

			s := Storage{db: db}

			err = s.DeleteToken(tt.givenID)
			if fmt.Sprint(err) != fmt.Sprint(tt.expectedErr) {
				t.Errorf("Returned error is not as expected. Expected:\n%q\nGiven:\n%q", tt.expectedErr, err)
			}
		})
	}
}
