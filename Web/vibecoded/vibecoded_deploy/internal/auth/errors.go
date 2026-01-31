package auth

import "errors"

// JWT errors
var (
	ErrMissingAuthCookie   = errors.New("token cookie is missing")
	ErrAuthCookieNotBearer = errors.New("token cookie is not Bearer")
	ErrInvalidToken        = errors.New("invalid token")
	ErrTokenFailedToParse  = errors.New("token failed to parse")
)

// Captcha errors
var (
	ErrMissingCaptchaCookie = errors.New("captcha cookie is missing")
	ErrWithCaptchaCookie    = errors.New("error with captcha cookie")
	ErrCaptchaTokenFailed   = errors.New("captcha token verification failed")
)
