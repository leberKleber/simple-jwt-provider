// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package web

import (
	"sync"
)

var (
	lockProviderMockCreatePasswordResetRequest sync.RWMutex
	lockProviderMockCreateUser                 sync.RWMutex
	lockProviderMockDeleteUser                 sync.RWMutex
	lockProviderMockLogin                      sync.RWMutex
	lockProviderMockResetPassword              sync.RWMutex
)

// Ensure, that ProviderMock does implement Provider.
// If this is not the case, regenerate this file with moq.
var _ Provider = &ProviderMock{}

// ProviderMock is a mock implementation of Provider.
//
//     func TestSomethingThatUsesProvider(t *testing.T) {
//
//         // make and configure a mocked Provider
//         mockedProvider := &ProviderMock{
//             CreatePasswordResetRequestFunc: func(email string) error {
// 	               panic("mock out the CreatePasswordResetRequest method")
//             },
//             CreateUserFunc: func(email string, password string, claims map[string]interface{}) error {
// 	               panic("mock out the CreateUser method")
//             },
//             DeleteUserFunc: func(email string) error {
// 	               panic("mock out the DeleteUser method")
//             },
//             LoginFunc: func(email string, password string) (string, error) {
// 	               panic("mock out the Login method")
//             },
//             ResetPasswordFunc: func(email string, resetToken string, password string) error {
// 	               panic("mock out the ResetPassword method")
//             },
//         }
//
//         // use mockedProvider in code that requires Provider
//         // and then make assertions.
//
//     }
type ProviderMock struct {
	// CreatePasswordResetRequestFunc mocks the CreatePasswordResetRequest method.
	CreatePasswordResetRequestFunc func(email string) error

	// CreateUserFunc mocks the CreateUser method.
	CreateUserFunc func(email string, password string, claims map[string]interface{}) error

	// DeleteUserFunc mocks the DeleteUser method.
	DeleteUserFunc func(email string) error

	// LoginFunc mocks the Login method.
	LoginFunc func(email string, password string) (string, error)

	// ResetPasswordFunc mocks the ResetPassword method.
	ResetPasswordFunc func(email string, resetToken string, password string) error

	// calls tracks calls to the methods.
	calls struct {
		// CreatePasswordResetRequest holds details about calls to the CreatePasswordResetRequest method.
		CreatePasswordResetRequest []struct {
			// Email is the email argument value.
			Email string
		}
		// CreateUser holds details about calls to the CreateUser method.
		CreateUser []struct {
			// Email is the email argument value.
			Email string
			// Password is the password argument value.
			Password string
			// Claims is the claims argument value.
			Claims map[string]interface{}
		}
		// DeleteUser holds details about calls to the DeleteUser method.
		DeleteUser []struct {
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
		// ResetPassword holds details about calls to the ResetPassword method.
		ResetPassword []struct {
			// Email is the email argument value.
			Email string
			// ResetToken is the resetToken argument value.
			ResetToken string
			// Password is the password argument value.
			Password string
		}
	}
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
	lockProviderMockCreatePasswordResetRequest.Lock()
	mock.calls.CreatePasswordResetRequest = append(mock.calls.CreatePasswordResetRequest, callInfo)
	lockProviderMockCreatePasswordResetRequest.Unlock()
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
	lockProviderMockCreatePasswordResetRequest.RLock()
	calls = mock.calls.CreatePasswordResetRequest
	lockProviderMockCreatePasswordResetRequest.RUnlock()
	return calls
}

// CreateUser calls CreateUserFunc.
func (mock *ProviderMock) CreateUser(email string, password string, claims map[string]interface{}) error {
	if mock.CreateUserFunc == nil {
		panic("ProviderMock.CreateUserFunc: method is nil but Provider.CreateUser was just called")
	}
	callInfo := struct {
		Email    string
		Password string
		Claims   map[string]interface{}
	}{
		Email:    email,
		Password: password,
		Claims:   claims,
	}
	lockProviderMockCreateUser.Lock()
	mock.calls.CreateUser = append(mock.calls.CreateUser, callInfo)
	lockProviderMockCreateUser.Unlock()
	return mock.CreateUserFunc(email, password, claims)
}

// CreateUserCalls gets all the calls that were made to CreateUser.
// Check the length with:
//     len(mockedProvider.CreateUserCalls())
func (mock *ProviderMock) CreateUserCalls() []struct {
	Email    string
	Password string
	Claims   map[string]interface{}
} {
	var calls []struct {
		Email    string
		Password string
		Claims   map[string]interface{}
	}
	lockProviderMockCreateUser.RLock()
	calls = mock.calls.CreateUser
	lockProviderMockCreateUser.RUnlock()
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
	lockProviderMockDeleteUser.Lock()
	mock.calls.DeleteUser = append(mock.calls.DeleteUser, callInfo)
	lockProviderMockDeleteUser.Unlock()
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
	lockProviderMockDeleteUser.RLock()
	calls = mock.calls.DeleteUser
	lockProviderMockDeleteUser.RUnlock()
	return calls
}

// Login calls LoginFunc.
func (mock *ProviderMock) Login(email string, password string) (string, error) {
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
	lockProviderMockLogin.Lock()
	mock.calls.Login = append(mock.calls.Login, callInfo)
	lockProviderMockLogin.Unlock()
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
	lockProviderMockLogin.RLock()
	calls = mock.calls.Login
	lockProviderMockLogin.RUnlock()
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
	lockProviderMockResetPassword.Lock()
	mock.calls.ResetPassword = append(mock.calls.ResetPassword, callInfo)
	lockProviderMockResetPassword.Unlock()
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
	lockProviderMockResetPassword.RLock()
	calls = mock.calls.ResetPassword
	lockProviderMockResetPassword.RUnlock()
	return calls
}
