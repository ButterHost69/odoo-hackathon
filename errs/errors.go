package errs

import "errors"

// AUTH
var (
	ErrInvalidCredentials       = errors.New("INVALID CREDENTIALS")
	ErrUsernameAlreadyUsedError = errors.New("EMAIL ALREADY USED")
	ErrSessionToken             = errors.New("INVALID SESSION TOKEN")
)

// USER
var (
	ErrUserEmailDoesNotExist = errors.New("USER EMAIL DOES NOT EXIST")
)

// AUTH
var (
	ErrSessionTokenDoesNotExist = errors.New("SESSION TOKEN DOES NOT EXIST")
)
