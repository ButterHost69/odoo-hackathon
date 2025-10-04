package errs

import "errors"

// External Errors(Display to Users)
const INTERNAL_SERVER_ERROR_MESSAGE = "<script>alert('Internal Server Error');</script>"

// Internal Errors
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
