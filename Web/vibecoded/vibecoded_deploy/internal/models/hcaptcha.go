package models

type HCaptchaVerifyResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes,omitempty"`
}

type HCaptchaVerifyRequest struct {
	Secret   string `json:"secret"`
	Response string `json:"response"`
}
