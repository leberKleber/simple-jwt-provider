// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package internal

import (
	"github.com/leberKleber/simple-jwt-provider/internal/storage"
	"sync"
)

// Ensure, that StorageMock does implement Storage.
// If this is not the case, regenerate this file with moq.
var _ Storage = &StorageMock{}

// StorageMock is a mock implementation of Storage.
//
// 	func TestSomethingThatUsesStorage(t *testing.T) {
//
// 		// make and configure a mocked Storage
// 		mockedStorage := &StorageMock{
// 			CreateTokenFunc: func(t *storage.Token) error {
// 				panic("mock out the CreateToken method")
// 			},
// 			CreateUserFunc: func(user storage.User) error {
// 				panic("mock out the CreateUser method")
// 			},
// 			DeleteTokenFunc: func(id uint) error {
// 				panic("mock out the DeleteToken method")
// 			},
// 			DeleteUserFunc: func(email string) error {
// 				panic("mock out the DeleteUser method")
// 			},
// 			TokensByEMailAndTokenFunc: func(email string, token string) ([]storage.Token, error) {
// 				panic("mock out the TokensByEMailAndToken method")
// 			},
// 			UpdateUserFunc: func(user storage.User) error {
// 				panic("mock out the UpdateUser method")
// 			},
// 			UserFunc: func(email string) (storage.User, error) {
// 				panic("mock out the User method")
// 			},
// 		}
//
// 		// use mockedStorage in code that requires Storage
// 		// and then make assertions.
//
// 	}
type StorageMock struct {
	// CreateTokenFunc mocks the CreateToken method.
	CreateTokenFunc func(t *storage.Token) error

	// CreateUserFunc mocks the CreateUser method.
	CreateUserFunc func(user storage.User) error

	// DeleteTokenFunc mocks the DeleteToken method.
	DeleteTokenFunc func(id uint) error

	// DeleteUserFunc mocks the DeleteUser method.
	DeleteUserFunc func(email string) error

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
			T *storage.Token
		}
		// CreateUser holds details about calls to the CreateUser method.
		CreateUser []struct {
			// User is the user argument value.
			User storage.User
		}
		// DeleteToken holds details about calls to the DeleteToken method.
		DeleteToken []struct {
			// ID is the id argument value.
			ID uint
		}
		// DeleteUser holds details about calls to the DeleteUser method.
		DeleteUser []struct {
			// Email is the email argument value.
			Email string
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
	lockCreateToken           sync.RWMutex
	lockCreateUser            sync.RWMutex
	lockDeleteToken           sync.RWMutex
	lockDeleteUser            sync.RWMutex
	lockTokensByEMailAndToken sync.RWMutex
	lockUpdateUser            sync.RWMutex
	lockUser                  sync.RWMutex
}

// CreateToken calls CreateTokenFunc.
func (mock *StorageMock) CreateToken(t *storage.Token) error {
	if mock.CreateTokenFunc == nil {
		panic("StorageMock.CreateTokenFunc: method is nil but Storage.CreateToken was just called")
	}
	callInfo := struct {
		T *storage.Token
	}{
		T: t,
	}
	mock.lockCreateToken.Lock()
	mock.calls.CreateToken = append(mock.calls.CreateToken, callInfo)
	mock.lockCreateToken.Unlock()
	return mock.CreateTokenFunc(t)
}

// CreateTokenCalls gets all the calls that were made to CreateToken.
// Check the length with:
//     len(mockedStorage.CreateTokenCalls())
func (mock *StorageMock) CreateTokenCalls() []struct {
	T *storage.Token
} {
	var calls []struct {
		T *storage.Token
	}
	mock.lockCreateToken.RLock()
	calls = mock.calls.CreateToken
	mock.lockCreateToken.RUnlock()
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
	mock.lockCreateUser.Lock()
	mock.calls.CreateUser = append(mock.calls.CreateUser, callInfo)
	mock.lockCreateUser.Unlock()
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
	mock.lockCreateUser.RLock()
	calls = mock.calls.CreateUser
	mock.lockCreateUser.RUnlock()
	return calls
}

// DeleteToken calls DeleteTokenFunc.
func (mock *StorageMock) DeleteToken(id uint) error {
	if mock.DeleteTokenFunc == nil {
		panic("StorageMock.DeleteTokenFunc: method is nil but Storage.DeleteToken was just called")
	}
	callInfo := struct {
		ID uint
	}{
		ID: id,
	}
	mock.lockDeleteToken.Lock()
	mock.calls.DeleteToken = append(mock.calls.DeleteToken, callInfo)
	mock.lockDeleteToken.Unlock()
	return mock.DeleteTokenFunc(id)
}

// DeleteTokenCalls gets all the calls that were made to DeleteToken.
// Check the length with:
//     len(mockedStorage.DeleteTokenCalls())
func (mock *StorageMock) DeleteTokenCalls() []struct {
	ID uint
} {
	var calls []struct {
		ID uint
	}
	mock.lockDeleteToken.RLock()
	calls = mock.calls.DeleteToken
	mock.lockDeleteToken.RUnlock()
	return calls
}

// DeleteUser calls DeleteUserFunc.
func (mock *StorageMock) DeleteUser(email string) error {
	if mock.DeleteUserFunc == nil {
		panic("StorageMock.DeleteUserFunc: method is nil but Storage.DeleteUser was just called")
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
//     len(mockedStorage.DeleteUserCalls())
func (mock *StorageMock) DeleteUserCalls() []struct {
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
	mock.lockTokensByEMailAndToken.Lock()
	mock.calls.TokensByEMailAndToken = append(mock.calls.TokensByEMailAndToken, callInfo)
	mock.lockTokensByEMailAndToken.Unlock()
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
	mock.lockTokensByEMailAndToken.RLock()
	calls = mock.calls.TokensByEMailAndToken
	mock.lockTokensByEMailAndToken.RUnlock()
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
	mock.lockUpdateUser.Lock()
	mock.calls.UpdateUser = append(mock.calls.UpdateUser, callInfo)
	mock.lockUpdateUser.Unlock()
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
	mock.lockUpdateUser.RLock()
	calls = mock.calls.UpdateUser
	mock.lockUpdateUser.RUnlock()
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
	mock.lockUser.Lock()
	mock.calls.User = append(mock.calls.User, callInfo)
	mock.lockUser.Unlock()
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
	mock.lockUser.RLock()
	calls = mock.calls.User
	mock.lockUser.RUnlock()
	return calls
}
