package errs

import "errors"

// External Errors(Display to Users)
const INTERNAL_SERVER_ERROR_MESSAGE = "<script>alert('Internal Server Error');</script>"
const UNAUTHORIZED_ACCESS_MESSAGE = "<script>alert('You're Not Authorized');</script>"

// Internal Errors
// AUTH
var (
	ErrSessionTokenDoesNotExist = errors.New("SESSION TOKEN DOES NOT EXIST")
	ErrInvalidCredentials       = errors.New("INVALID CREDENTIALS")
	ErrUsernameAlreadyUsedError = errors.New("EMAIL ALREADY USED")
	ErrSessionToken             = errors.New("INVALID SESSION TOKEN")
)

// COMPANY
var (
	ErrAdminEmailNotFound = errors.New("COMPANY EMAIL NOT FOUND")
)

// USER
var (
	ErrUserEmailDoesNotExist = errors.New("USER EMAIL DOES NOT EXIST")
)

// UTILS
var ErrCountryNotFound = errors.New("country not found")
var ErrCurrencyNotFound = errors.New("no currency found")
