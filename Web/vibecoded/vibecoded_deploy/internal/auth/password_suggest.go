package auth

import (
	"math/rand"
)

const (
	MinLength      = 8
	MaxLength      = 32
	digitChars     = "0123456789"
	lowercaseChars = "abcdefghijklmnopqrstuvwxyz"
	uppercaseChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialChars   = `!"#$%&\'()*+,-./:;<=>?@[\]^_{|}~` + "`"
	allChars       = digitChars + lowercaseChars + uppercaseChars + specialChars
)

var CharGroups = map[string]string{
	"digit":     digitChars,
	"lowercase": lowercaseChars,
	"uppercase": uppercaseChars,
	"special":   specialChars,
}

func addRandomChar(password, chars string) string {
	randomChar := string(chars[rand.Intn(len(chars))])
	if rand.Intn(2) == 0 {
		return password + randomChar
	}
	return randomChar + password
}

func generateSuggestion(initialPassword string) string {
	maxInitialPasswordLength := MaxLength - 2
	if len(initialPassword) > maxInitialPasswordLength {
		initialPassword = initialPassword[:maxInitialPasswordLength]
	}

	for range 2 {
		initialPassword = addRandomChar(initialPassword, allChars)
	}
	return initialPassword
}

func suggestPasswords(initialPassword string) []string {
	var suggestedPasswords []string
	for range 3 {
		suggestedPasswords = append(suggestedPasswords, generateSuggestion(initialPassword))
	}
	return suggestedPasswords
}
