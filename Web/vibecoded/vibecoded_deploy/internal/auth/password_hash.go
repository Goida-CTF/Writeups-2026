package auth

import (
	"crypto/sha256"
	"encoding/hex"
)

type PasswordHasher struct {
	pepper []byte
}

func NewPasswordHasher(pepper []byte) *PasswordHasher {
	return &PasswordHasher{
		pepper: pepper,
	}
}

func (ph *PasswordHasher) HashPassword(password string) string {
	data := append([]byte{}, ph.pepper...)
	data = append(data, []byte(password)...)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
