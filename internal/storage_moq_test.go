// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package internal

import (
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"sync"
)

var (
	lockStorageMockCreateToken           sync.RWMutex
	lockStorageMockCreateUser            sync.RWMutex
	lockStorageMockDeleteToken           sync.RWMutex
	lockStorageMockTokensByEMailAndToken sync.RWMutex
	lockStorageMockUpdateUser            sync.RWMutex
	lockStorageMockUser                  sync.RWMutex
)

// Ensure, that StorageMock does implement Storage.
// If this is not the case, regenerate this file with moq.
var _ Storage = &StorageMock{}

// StorageMock is a mock implementation of Storage.
//
//     func TestSomethingThatUsesStorage(t *testing.T) {
//
//         // make and configure a mocked Storage
//         mockedStorage := &StorageMock{
//             CreateTokenFunc: func(t storage.Token) (int64, error) {
// 	               panic("mock out the CreateToken method")
//             },
//             CreateUserFunc: func(user storage.User) error {
// 	               panic("mock out the CreateUser method")
//             },
//             DeleteTokenFunc: func(id int64) error {
// 	               panic("mock out the DeleteToken method")
//             },
//             TokensByEMailAndTokenFunc: func(email string, token string) ([]storage.Token, error) {
// 	               panic("mock out the TokensByEMailAndToken method")
//             },
//             UpdateUserFunc: func(user storage.User) error {
// 	               panic("mock out the UpdateUser method")
//             },
//             UserFunc: func(email string) (storage.User, error) {
// 	               panic("mock out the User method")
//             },
//         }
//
//         // use mockedStorage in code that requires Storage
//         // and then make assertions.
//
//     }
type StorageMock struct {
	// CreateTokenFunc mocks the CreateToken method.
	CreateTokenFunc func(t storage.Token) (int64, error)

	// CreateUserFunc mocks the CreateUser method.
	CreateUserFunc func(user storage.User) error

	// DeleteTokenFunc mocks the DeleteToken method.
	DeleteTokenFunc func(id int64) error

	// TokensByEMailAndTokenFunc mocks the TokensByEMailAndToken method.
	TokensByEMailAndTokenFunc func(email string, token string) ([]storage.Token, error)

	// UpdateUserFunc mocks the UpdateUser method.
	UpdateUserFunc func(user storage.User) error

	// UserFunc mocks the User method.
	UserFunc func(email string) (storage.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// CreateToken holds details about calls to the CreateToken method.
		CreateToken []struct {
			// T is the t argument value.
			T storage.Token
		}
		// CreateUser holds details about calls to the CreateUser method.
		CreateUser []struct {
			// User is the user argument value.
			User storage.User
		}
		// DeleteToken holds details about calls to the DeleteToken method.
		DeleteToken []struct {
			// ID is the id argument value.
			ID int64
		}
		// TokensByEMailAndToken holds details about calls to the TokensByEMailAndToken method.
		TokensByEMailAndToken []struct {
			// Email is the email argument value.
			Email string
			// Token is the token argument value.
			Token string
		}
		// UpdateUser holds details about calls to the UpdateUser method.
		UpdateUser []struct {
			// User is the user argument value.
			User storage.User
		}
		// User holds details about calls to the User method.
		User []struct {
			// Email is the email argument value.
			Email string
		}
	}
}

// CreateToken calls CreateTokenFunc.
func (mock *StorageMock) CreateToken(t storage.Token) (int64, error) {
	if mock.CreateTokenFunc == nil {
		panic("StorageMock.CreateTokenFunc: method is nil but Storage.CreateToken was just called")
	}
	callInfo := struct {
		T storage.Token
	}{
		T: t,
	}
	lockStorageMockCreateToken.Lock()
	mock.calls.CreateToken = append(mock.calls.CreateToken, callInfo)
	lockStorageMockCreateToken.Unlock()
	return mock.CreateTokenFunc(t)
}

// CreateTokenCalls gets all the calls that were made to CreateToken.
// Check the length with:
//     len(mockedStorage.CreateTokenCalls())
func (mock *StorageMock) CreateTokenCalls() []struct {
	T storage.Token
} {
	var calls []struct {
		T storage.Token
	}
	lockStorageMockCreateToken.RLock()
	calls = mock.calls.CreateToken
	lockStorageMockCreateToken.RUnlock()
	return calls
}

// CreateUser calls CreateUserFunc.
func (mock *StorageMock) CreateUser(user storage.User) error {
	if mock.CreateUserFunc == nil {
		panic("StorageMock.CreateUserFunc: method is nil but Storage.CreateUser was just called")
	}
	callInfo := struct {
		User storage.User
	}{
		User: user,
	}
	lockStorageMockCreateUser.Lock()
	mock.calls.CreateUser = append(mock.calls.CreateUser, callInfo)
	lockStorageMockCreateUser.Unlock()
	return mock.CreateUserFunc(user)
}

// CreateUserCalls gets all the calls that were made to CreateUser.
// Check the length with:
//     len(mockedStorage.CreateUserCalls())
func (mock *StorageMock) CreateUserCalls() []struct {
	User storage.User
} {
	var calls []struct {
		User storage.User
	}
	lockStorageMockCreateUser.RLock()
	calls = mock.calls.CreateUser
	lockStorageMockCreateUser.RUnlock()
	return calls
}

// DeleteToken calls DeleteTokenFunc.
func (mock *StorageMock) DeleteToken(id int64) error {
	if mock.DeleteTokenFunc == nil {
		panic("StorageMock.DeleteTokenFunc: method is nil but Storage.DeleteToken was just called")
	}
	callInfo := struct {
		ID int64
	}{
		ID: id,
	}
	lockStorageMockDeleteToken.Lock()
	mock.calls.DeleteToken = append(mock.calls.DeleteToken, callInfo)
	lockStorageMockDeleteToken.Unlock()
	return mock.DeleteTokenFunc(id)
}

// DeleteTokenCalls gets all the calls that were made to DeleteToken.
// Check the length with:
//     len(mockedStorage.DeleteTokenCalls())
func (mock *StorageMock) DeleteTokenCalls() []struct {
	ID int64
} {
	var calls []struct {
		ID int64
	}
	lockStorageMockDeleteToken.RLock()
	calls = mock.calls.DeleteToken
	lockStorageMockDeleteToken.RUnlock()
	return calls
}

// TokensByEMailAndToken calls TokensByEMailAndTokenFunc.
func (mock *StorageMock) TokensByEMailAndToken(email string, token string) ([]storage.Token, error) {
	if mock.TokensByEMailAndTokenFunc == nil {
		panic("StorageMock.TokensByEMailAndTokenFunc: method is nil but Storage.TokensByEMailAndToken was just called")
	}
	callInfo := struct {
		Email string
		Token string
	}{
		Email: email,
		Token: token,
	}
	lockStorageMockTokensByEMailAndToken.Lock()
	mock.calls.TokensByEMailAndToken = append(mock.calls.TokensByEMailAndToken, callInfo)
	lockStorageMockTokensByEMailAndToken.Unlock()
	return mock.TokensByEMailAndTokenFunc(email, token)
}

// TokensByEMailAndTokenCalls gets all the calls that were made to TokensByEMailAndToken.
// Check the length with:
//     len(mockedStorage.TokensByEMailAndTokenCalls())
func (mock *StorageMock) TokensByEMailAndTokenCalls() []struct {
	Email string
	Token string
} {
	var calls []struct {
		Email string
		Token string
	}
	lockStorageMockTokensByEMailAndToken.RLock()
	calls = mock.calls.TokensByEMailAndToken
	lockStorageMockTokensByEMailAndToken.RUnlock()
	return calls
}

// UpdateUser calls UpdateUserFunc.
func (mock *StorageMock) UpdateUser(user storage.User) error {
	if mock.UpdateUserFunc == nil {
		panic("StorageMock.UpdateUserFunc: method is nil but Storage.UpdateUser was just called")
	}
	callInfo := struct {
		User storage.User
	}{
		User: user,
	}
	lockStorageMockUpdateUser.Lock()
	mock.calls.UpdateUser = append(mock.calls.UpdateUser, callInfo)
	lockStorageMockUpdateUser.Unlock()
	return mock.UpdateUserFunc(user)
}

// UpdateUserCalls gets all the calls that were made to UpdateUser.
// Check the length with:
//     len(mockedStorage.UpdateUserCalls())
func (mock *StorageMock) UpdateUserCalls() []struct {
	User storage.User
} {
	var calls []struct {
		User storage.User
	}
	lockStorageMockUpdateUser.RLock()
	calls = mock.calls.UpdateUser
	lockStorageMockUpdateUser.RUnlock()
	return calls
}

// User calls UserFunc.
func (mock *StorageMock) User(email string) (storage.User, error) {
	if mock.UserFunc == nil {
		panic("StorageMock.UserFunc: method is nil but Storage.User was just called")
	}
	callInfo := struct {
		Email string
	}{
		Email: email,
	}
	lockStorageMockUser.Lock()
	mock.calls.User = append(mock.calls.User, callInfo)
	lockStorageMockUser.Unlock()
	return mock.UserFunc(email)
}

// UserCalls gets all the calls that were made to User.
// Check the length with:
//     len(mockedStorage.UserCalls())
func (mock *StorageMock) UserCalls() []struct {
	Email string
} {
	var calls []struct {
		Email string
	}
	lockStorageMockUser.RLock()
	calls = mock.calls.User
	lockStorageMockUser.RUnlock()
	return calls
}