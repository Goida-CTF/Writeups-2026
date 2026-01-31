package auth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"

	"vibecoded/internal/models"
)

const (
	hCaptchaVerifyURL = "https://api.hcaptcha.com/siteverify"
	CaptchaCookieName = "_c"
)

type CaptchaProvider struct {
	secret string
}

func NewCaptchaProvider(captchaSecret string) *CaptchaProvider {
	return &CaptchaProvider{
		secret: captchaSecret,
	}
}

func (c *CaptchaProvider) VerifyCaptchaToken(token string) (bool, error) {
	resp, err := http.PostForm(hCaptchaVerifyURL, url.Values{
		"secret":   {c.secret},
		"response": {token},
	})
	if err != nil {
		return false, fmt.Errorf("http.Post: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("io.ReadAll: %w", err)
	}

	var hCaptchaResp models.HCaptchaVerifyResponse
	if err := json.Unmarshal(body, &hCaptchaResp); err != nil {
		return false, fmt.Errorf("json.Unmarshal: %w", err)
	}

	log.Debugf("hCaptcha was verified with response: %s", string(body))
	return hCaptchaResp.Success, nil
}

func (c *CaptchaProvider) HandleVerification(r *http.Request) (bool, error) {
	tokenCookie, err := r.Cookie(CaptchaCookieName)
	if err != nil || tokenCookie.Value == "" {
		return false, ErrMissingCaptchaCookie
	}

	ok, err := c.VerifyCaptchaToken(tokenCookie.Value)
	if err != nil {
		return false, fmt.Errorf("c.VerifyCaptchaToken: %w", err)
	}

	return ok, nil
}
