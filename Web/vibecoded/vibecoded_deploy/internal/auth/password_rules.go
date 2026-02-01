package auth

import (
	"fmt"
	"strings"

	"vibecoded/internal/models"
)

const (
	textTooShort    = "password is too short, minimum 8 characters"
	textNoDigits    = "password should contain at least one digit"
	textNoUppercase = "password should contain at least one uppercase letter"
	textNoLowercase = "password should contain at least one lowercase letter"
	textNoSpecial   = "password should contain at least one special character"
	textTooLong     = "password is too secure: maximum 32 characters"
	textUnallowed   = "password is too secure: contains unallowed characters"
	textCollisionF  = "password is already taken by user %s\nsuggested passwords: %s"
)

var charToTextMap = map[string]string{
	"digit":     textNoDigits,
	"lowercase": textNoLowercase,
	"uppercase": textNoUppercase,
	"special":   textNoSpecial,
}

func verifyPasswordRules(password string) (int, string) {
	var level = 3

	for _, char := range password {
		if !strings.ContainsRune(allChars, char) {
			return level, textUnallowed
		}
	}
	if len(password) > MaxLength {
		return level, textTooLong
	}

	level = 1

	if len(password) < MinLength {
		return level, textTooShort
	}

	level = 2

	for charGroupName, charGroupChars := range CharGroups {
		if !strings.ContainsAny(password, charGroupChars) {
			return level, charToTextMap[charGroupName]
		}
	}

	return level, ""
}

func CheckPasswordComplexity(password, collisionUsername string,
) *models.PasswordComplexityResult {
	var (
		level = 1
		text  string
	)

	if collisionUsername != "" {
		suggestedPasswords := strings.Join(suggestPasswords(password), " ")

		return &models.PasswordComplexityResult{
			Ok:    false,
			Level: level,
			Text: fmt.Sprintf(textCollisionF,
				collisionUsername, suggestedPasswords),
		}
	}

	level, text = verifyPasswordRules(password)

	return &models.PasswordComplexityResult{
		Ok:    text == "",
		Level: level,
		Text:  text,
	}
}
