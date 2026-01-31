package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"

	"vibecoded/internal/auth"
	"vibecoded/internal/models"
)

const InternalServerError = "Internal Server Error"

var (
	ErrEmptyRequestBody   = errors.New("empty request body")
	ErrEmptyField         = errors.New("field is empty")
	ErrContentTypeNotJSON = errors.New("expect application/json in Content-Type")
)

func parseJWTError(err error) error {
	for {
		e := errors.Unwrap(err)
		if e == nil {
			break
		}
		err = e
	}

	if !(errors.Is(err, jwt.ErrInvalidKey) ||
		errors.Is(err, jwt.ErrTokenExpired) ||
		errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenSignatureInvalid) ||
		errors.Is(err, jwt.ErrTokenUnverifiable) ||
		errors.Is(err, auth.ErrMissingAuthCookie) ||
		errors.Is(err, auth.ErrAuthCookieNotBearer) ||
		errors.Is(err, auth.ErrInvalidToken)) {
		err = auth.ErrTokenFailedToParse
	}
	return err
}

func handleJWTError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithError(err).Debugln("failed to parse JWT token")

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Unix(0, 0),
		Path:     "/",
		HttpOnly: true,
	})

	parsedError := parseJWTError(err)
	HandleClientError(w, r,
		err,
		parsedError.Error(),
		http.StatusUnauthorized)
}

func handleCaptchaError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithError(err).Debugln("failed to verify Captcha")

	http.SetCookie(w, &http.Cookie{
		Name:    auth.CaptchaCookieName,
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	var errForMsg error
	if errors.Is(err, auth.ErrMissingCaptchaCookie) {
		errForMsg = err
	} else {
		errForMsg = auth.ErrWithCaptchaCookie
	}

	HandleClientError(w, r,
		err,
		errForMsg.Error(),
		http.StatusBadRequest)
}

func HandleInternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithError(err).Errorf("%s %s - %s", r.Method, r.URL.Path, InternalServerError)

	res := &models.JsonResponse{
		Status: models.ResponseStatusError,
		Error:  InternalServerError,
	}

	if err := renderJSONWithCode(w, http.StatusInternalServerError, &res); err != nil {
		log.WithError(err).Errorf("%s %s - %s", r.Method, r.URL.Path, InternalServerError)
		http.Error(w, InternalServerError, http.StatusInternalServerError)
	}
}

func HandleClientError(w http.ResponseWriter, r *http.Request, err error, msg string, code int) {
	log.WithError(err).Debugln(msg)

	res := &models.JsonResponse{
		Status: models.ResponseStatusError,
		Error:  msg,
	}

	if err := renderJSONWithCode(w, code, &res); err != nil {
		log.WithError(err).Errorf("HandleClientError: %s", msg)
		http.Error(w, msg, code)
	}
}
