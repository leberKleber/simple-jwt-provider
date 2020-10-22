// +build component

package main

import (
	"fmt"
	"testing"
)

func TestCRUDUser(t *testing.T) {
	// 1) create user
	// 2) login
	// 3) update user
	// 4) login
	// 5) get user
	// 6) delete user
	// 7) login

	email := "crud_test@leberkleber.io"
	password := "s3cr3t"

	// 1)
	createUser(t, email, password)

	// 2)
	loginUser(t, email, password)

	// 3)
	newPassword := "n3wS3cr3t"
	updateUser(t, email, newPassword, map[string]interface{}{
		"myClaim": 5,
	})

	// 4)
	loginUser(t, email, newPassword)

	// 5)
	expectedUser := User{
		EMail:    email,
		Password: "**********",
		Claims: map[string]interface{}{
			"myClaim": 5,
		},
	}
	user := readUser(t, email)
	if fmt.Sprint(user) != fmt.Sprint(expectedUser) {
		t.Fatalf("user is not as expected. Expected:\n%#v\nGiven:\n%#v", expectedUser, user)
	}

	// 6)
	deleteUser(t, email)

	// 7)
	loginUser(t, email, newPassword)
}
