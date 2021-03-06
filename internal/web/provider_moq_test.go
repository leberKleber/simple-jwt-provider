// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package web

import (
	"github.com/leberKleber/simple-jwt-provider/internal"
	"sync"
)

// Ensure, that ProviderMock does implement Provider.
// If this is not the case, regenerate this file with moq.
var _ Provider = &ProviderMock{}

// ProviderMock is a mock implementation of Provider.
//
// 	func TestSomethingThatUsesProvider(t *testing.T) {
//
// 		// make and configure a mocked Provider
// 		mockedProvider := &ProviderMock{
// 			CreatePasswordResetRequestFunc: func(email string) error {
// 				panic("mock out the CreatePasswordResetRequest method")
// 			},
// 			CreateUserFunc: func(user internal.User) error {
// 				panic("mock out the CreateUser method")
// 			},
// 			DeleteUserFunc: func(email string) error {
// 				panic("mock out the DeleteUser method")
// 			},
// 			GetUserFunc: func(email string) (internal.User, error) {
// 				panic("mock out the GetUser method")
// 			},
// 			LoginFunc: func(email string, password string) (string, string, error) {
// 				panic("mock out the Login method")
// 			},
// 			RefreshFunc: func(refreshToken string) (string, string, error) {
// 				panic("mock out the Refresh method")
// 			},
// 			ResetPasswordFunc: func(email string, resetToken string, password string) error {
// 				panic("mock out the ResetPassword method")
// 			},
// 			UpdateUserFunc: func(email string, user internal.User) (internal.User, error) {
// 				panic("mock out the UpdateUser method")
// 			},
// 		}
//
// 		// use mockedProvider in code that requires Provider
// 		// and then make assertions.
//
// 	}
type ProviderMock struct {
	// CreatePasswordResetRequestFunc mocks the CreatePasswordResetRequest method.
	CreatePasswordResetRequestFunc func(email string) error

	// CreateUserFunc mocks the CreateUser method.
	CreateUserFunc func(user internal.User) error

	// DeleteUserFunc mocks the DeleteUser method.
	DeleteUserFunc func(email string) error

	// GetUserFunc mocks the GetUser method.
	GetUserFunc func(email string) (internal.User, error)

	// LoginFunc mocks the Login method.
	LoginFunc func(email string, password string) (string, string, error)

	// RefreshFunc mocks the Refresh method.
	RefreshFunc func(refreshToken string) (string, string, error)

	// ResetPasswordFunc mocks the ResetPassword method.
	ResetPasswordFunc func(email string, resetToken string, password string) error

	// UpdateUserFunc mocks the UpdateUser method.
	UpdateUserFunc func(email string, user internal.User) (internal.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreatePasswordResetRequest holds details about calls to the CreatePasswordResetRequest method.
		CreatePasswordResetRequest []struct {
			// Email is the email argument value.
			Email string
		}
		// CreateUser holds details about calls to the CreateUser method.
		CreateUser []struct {
			// User is the user argument value.
			User internal.User
		}
		// DeleteUser holds details about calls to the DeleteUser method.
		DeleteUser []struct {
			// Email is the email argument value.
			Email string
		}
		// GetUser holds details about calls to the GetUser method.
		GetUser []struct {
			// Email is the email argument value.
			Email string
		}
		// Login holds details about calls to the Login method.
		Login []struct {
			// Email is the email argument value.
			Email string
			// Password is the password argument value.
			Password string
		}
		// Refresh holds details about calls to the Refresh method.
		Refresh []struct {
			// RefreshToken is the refreshToken argument value.
			RefreshToken string
		}
		// ResetPassword holds details about calls to the ResetPassword method.
		ResetPassword []struct {
			// Email is the email argument value.
			Email string
			// ResetToken is the resetToken argument value.
			ResetToken string
			// Password is the password argument value.
			Password string
		}
		// UpdateUser holds details about calls to the UpdateUser method.
		UpdateUser []struct {
			// Email is the email argument value.
			Email string
			// User is the user argument value.
			User internal.User
		}
	}
	lockCreatePasswordResetRequest sync.RWMutex
	lockCreateUser                 sync.RWMutex
	lockDeleteUser                 sync.RWMutex
	lockGetUser                    sync.RWMutex
	lockLogin                      sync.RWMutex
	lockRefresh                    sync.RWMutex
	lockResetPassword              sync.RWMutex
	lockUpdateUser                 sync.RWMutex
}

// CreatePasswordResetRequest calls CreatePasswordResetRequestFunc.
func (mock *ProviderMock) CreatePasswordResetRequest(email string) error {
	if mock.CreatePasswordResetRequestFunc == nil {
		panic("ProviderMock.CreatePasswordResetRequestFunc: method is nil but Provider.CreatePasswordResetRequest was just called")
	}
	callInfo := struct {
		Email string
	}{
		Email: email,
	}
	mock.lockCreatePasswordResetRequest.Lock()
	mock.calls.CreatePasswordResetRequest = append(mock.calls.CreatePasswordResetRequest, callInfo)
	mock.lockCreatePasswordResetRequest.Unlock()
	return mock.CreatePasswordResetRequestFunc(email)
}

// CreatePasswordResetRequestCalls gets all the calls that were made to CreatePasswordResetRequest.
// Check the length with:
//     len(mockedProvider.CreatePasswordResetRequestCalls())
func (mock *ProviderMock) CreatePasswordResetRequestCalls() []struct {
	Email string
} {
	var calls []struct {
		Email string
	}
	mock.lockCreatePasswordResetRequest.RLock()
	calls = mock.calls.CreatePasswordResetRequest
	mock.lockCreatePasswordResetRequest.RUnlock()
	return calls
}

// CreateUser calls CreateUserFunc.
func (mock *ProviderMock) CreateUser(user internal.User) error {
	if mock.CreateUserFunc == nil {
		panic("ProviderMock.CreateUserFunc: method is nil but Provider.CreateUser was just called")
	}
	callInfo := struct {
		User internal.User
	}{
		User: user,
	}
	mock.lockCreateUser.Lock()
	mock.calls.CreateUser = append(mock.calls.CreateUser, callInfo)
	mock.lockCreateUser.Unlock()
	return mock.CreateUserFunc(user)
}

// CreateUserCalls gets all the calls that were made to CreateUser.
// Check the length with:
//     len(mockedProvider.CreateUserCalls())
func (mock *ProviderMock) CreateUserCalls() []struct {
	User internal.User
} {
	var calls []struct {
		User internal.User
	}
	mock.lockCreateUser.RLock()
	calls = mock.calls.CreateUser
	mock.lockCreateUser.RUnlock()
	return calls
}

// DeleteUser calls DeleteUserFunc.
func (mock *ProviderMock) DeleteUser(email string) error {
	if mock.DeleteUserFunc == nil {
		panic("ProviderMock.DeleteUserFunc: method is nil but Provider.DeleteUser was just called")
	}
	callInfo := struct {
		Email string
	}{
		Email: email,
	}
	mock.lockDeleteUser.Lock()
	mock.calls.DeleteUser = append(mock.calls.DeleteUser, callInfo)
	mock.lockDeleteUser.Unlock()
	return mock.DeleteUserFunc(email)
}

// DeleteUserCalls gets all the calls that were made to DeleteUser.
// Check the length with:
//     len(mockedProvider.DeleteUserCalls())
func (mock *ProviderMock) DeleteUserCalls() []struct {
	Email string
} {
	var calls []struct {
		Email string
	}
	mock.lockDeleteUser.RLock()
	calls = mock.calls.DeleteUser
	mock.lockDeleteUser.RUnlock()
	return calls
}

// GetUser calls GetUserFunc.
func (mock *ProviderMock) GetUser(email string) (internal.User, error) {
	if mock.GetUserFunc == nil {
		panic("ProviderMock.GetUserFunc: method is nil but Provider.GetUser was just called")
	}
	callInfo := struct {
		Email string
	}{
		Email: email,
	}
	mock.lockGetUser.Lock()
	mock.calls.GetUser = append(mock.calls.GetUser, callInfo)
	mock.lockGetUser.Unlock()
	return mock.GetUserFunc(email)
}

// GetUserCalls gets all the calls that were made to GetUser.
// Check the length with:
//     len(mockedProvider.GetUserCalls())
func (mock *ProviderMock) GetUserCalls() []struct {
	Email string
} {
	var calls []struct {
		Email string
	}
	mock.lockGetUser.RLock()
	calls = mock.calls.GetUser
	mock.lockGetUser.RUnlock()
	return calls
}

// Login calls LoginFunc.
func (mock *ProviderMock) Login(email string, password string) (string, string, error) {
	if mock.LoginFunc == nil {
		panic("ProviderMock.LoginFunc: method is nil but Provider.Login was just called")
	}
	callInfo := struct {
		Email    string
		Password string
	}{
		Email:    email,
		Password: password,
	}
	mock.lockLogin.Lock()
	mock.calls.Login = append(mock.calls.Login, callInfo)
	mock.lockLogin.Unlock()
	return mock.LoginFunc(email, password)
}

// LoginCalls gets all the calls that were made to Login.
// Check the length with:
//     len(mockedProvider.LoginCalls())
func (mock *ProviderMock) LoginCalls() []struct {
	Email    string
	Password string
} {
	var calls []struct {
		Email    string
		Password string
	}
	mock.lockLogin.RLock()
	calls = mock.calls.Login
	mock.lockLogin.RUnlock()
	return calls
}

// Refresh calls RefreshFunc.
func (mock *ProviderMock) Refresh(refreshToken string) (string, string, error) {
	if mock.RefreshFunc == nil {
		panic("ProviderMock.RefreshFunc: method is nil but Provider.Refresh was just called")
	}
	callInfo := struct {
		RefreshToken string
	}{
		RefreshToken: refreshToken,
	}
	mock.lockRefresh.Lock()
	mock.calls.Refresh = append(mock.calls.Refresh, callInfo)
	mock.lockRefresh.Unlock()
	return mock.RefreshFunc(refreshToken)
}

// RefreshCalls gets all the calls that were made to Refresh.
// Check the length with:
//     len(mockedProvider.RefreshCalls())
func (mock *ProviderMock) RefreshCalls() []struct {
	RefreshToken string
} {
	var calls []struct {
		RefreshToken string
	}
	mock.lockRefresh.RLock()
	calls = mock.calls.Refresh
	mock.lockRefresh.RUnlock()
	return calls
}

// ResetPassword calls ResetPasswordFunc.
func (mock *ProviderMock) ResetPassword(email string, resetToken string, password string) error {
	if mock.ResetPasswordFunc == nil {
		panic("ProviderMock.ResetPasswordFunc: method is nil but Provider.ResetPassword was just called")
	}
	callInfo := struct {
		Email      string
		ResetToken string
		Password   string
	}{
		Email:      email,
		ResetToken: resetToken,
		Password:   password,
	}
	mock.lockResetPassword.Lock()
	mock.calls.ResetPassword = append(mock.calls.ResetPassword, callInfo)
	mock.lockResetPassword.Unlock()
	return mock.ResetPasswordFunc(email, resetToken, password)
}

// ResetPasswordCalls gets all the calls that were made to ResetPassword.
// Check the length with:
//     len(mockedProvider.ResetPasswordCalls())
func (mock *ProviderMock) ResetPasswordCalls() []struct {
	Email      string
	ResetToken string
	Password   string
} {
	var calls []struct {
		Email      string
		ResetToken string
		Password   string
	}
	mock.lockResetPassword.RLock()
	calls = mock.calls.ResetPassword
	mock.lockResetPassword.RUnlock()
	return calls
}

// UpdateUser calls UpdateUserFunc.
func (mock *ProviderMock) UpdateUser(email string, user internal.User) (internal.User, error) {
	if mock.UpdateUserFunc == nil {
		panic("ProviderMock.UpdateUserFunc: method is nil but Provider.UpdateUser was just called")
	}
	callInfo := struct {
		Email string
		User  internal.User
	}{
		Email: email,
		User:  user,
	}
	mock.lockUpdateUser.Lock()
	mock.calls.UpdateUser = append(mock.calls.UpdateUser, callInfo)
	mock.lockUpdateUser.Unlock()
	return mock.UpdateUserFunc(email, user)
}

// UpdateUserCalls gets all the calls that were made to UpdateUser.
// Check the length with:
//     len(mockedProvider.UpdateUserCalls())
func (mock *ProviderMock) UpdateUserCalls() []struct {
	Email string
	User  internal.User
} {
	var calls []struct {
		Email string
		User  internal.User
	}
	mock.lockUpdateUser.RLock()
	calls = mock.calls.UpdateUser
	mock.lockUpdateUser.RUnlock()
	return calls
}
